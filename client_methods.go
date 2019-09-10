package makerbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/tjhorner/makerbot-rpc/printfile"
	"github.com/tjhorner/makerbot-rpc/reflector"
)

func (c *Client) ping() (*bool, error) {
	var reply bool
	return &reply, c.call("ping", rpcEmptyParams{}, &reply)
}

func (c *Client) sendHandshake() (*Printer, error) {
	var reply Printer
	return &reply, c.call("handshake", rpcEmptyParams{}, &reply)
}

type rpcAuthenticateParams struct {
	AccessToken string `json:"access_token"`
}

// authenticate performs authentication with the printer
// via an access token retrieved through the printer's
// HTTP server
func (c *Client) authenticate(accessToken string) (*json.RawMessage, error) {
	var reply json.RawMessage
	return &reply, c.call("authenticate", rpcAuthenticateParams{accessToken}, &reply)
}

type rpcAuthPacketParams struct {
	CallID     string `json:"call_id"`
	ClientCode string `json:"client_code"`
	PrinterID  string `json:"printer_id"`
}

func (c *Client) sendAuthPacket(id string, pc *reflector.CallPrinterResponse) (*bool, error) {
	params := rpcAuthPacketParams{
		CallID:     pc.Call.ID,
		ClientCode: pc.Call.ClientCode,
		PrinterID:  id,
	}

	var reply bool
	return &reply, c.call("auth_packet", params, &reply)
}

type rpcLoadUnloadFilamentParams struct {
	ToolIndex int `json:"tool_index"`
}

// LoadFilament instructs the printer to begin loading filament into
// the extruder
func (c *Client) LoadFilament(toolIndex int) (*PrinterProcess, error) {
	var reply PrinterProcess
	return &reply, c.call("load_filament", rpcLoadUnloadFilamentParams{toolIndex}, &reply)
}

// UnloadFilament instructs the printer to begin unloading filament from
// the extruder
func (c *Client) UnloadFilament(toolIndex int) (*json.RawMessage, error) {
	var reply json.RawMessage
	return &reply, c.call("unload_filament", rpcLoadUnloadFilamentParams{toolIndex}, &reply)
}

// Cancel instructs the printer to cancel the current process, if any.
//
// It may result in a `ProcessNotCancellableException`, so you may want to
// check the `CurrentProcess` to ensure it is `Cancellable`. Or not, if you
// don't care if it fails.
func (c *Client) Cancel() (*json.RawMessage, error) {
	var reply json.RawMessage
	return &reply, c.call("cancel", rpcEmptyParams{}, &reply)
}

type rpcProcessMethodParams struct {
	Method string `json:"method"`
}

// ProcessMethod will send a process_method request to the printer with no parameters.
func (c *Client) ProcessMethod(method string) (*json.RawMessage, error) {
	var reply json.RawMessage
	return &reply, c.call("process_method", rpcProcessMethodParams{method}, &reply)
}

// Suspend instructs the printer to suspend the current process, if any.
//
// Suspend can be reversed by using Resume.
func (c *Client) Suspend() (*json.RawMessage, error) {
	var reply json.RawMessage
	return &reply, c.call("process_method", rpcProcessMethodParams{"suspend"}, &reply)
}

// Resume instructs the printer to resume the current process, if any.
//
// Resume can be reversed by using Suspend.
func (c *Client) Resume() (*json.RawMessage, error) {
	var reply json.RawMessage
	return &reply, c.call("process_method", rpcProcessMethodParams{"resume"}, &reply)
}

type rpcChangeMachineNameParams struct {
	MachineName string `json:"machine_name"`
}

// ChangeMachineName instructs the printer to change its display name.
func (c *Client) ChangeMachineName(name string) (*json.RawMessage, error) {
	var reply json.RawMessage
	return &reply, c.call("cancel", rpcChangeMachineNameParams{name}, &reply)
}

func (c *Client) requestCameraFrame() (*bool, error) {
	var reply bool
	return &reply, c.call("request_camera_frame", rpcEmptyParams{}, &reply)
}

func (c *Client) requestCameraStream() error {
	return c.call("request_camera_stream", rpcEmptyParams{}, nil)
}

func (c *Client) endCameraStream() error {
	return c.call("end_camera_stream", rpcEmptyParams{}, nil)
}

// GetCameraFrame requests a single frame from the printer's camera
func (c *Client) GetCameraFrame() (*CameraFrame, error) {
	ch := make(chan CameraFrame)
	c.cameraCh = &ch

	res, err := c.requestCameraFrame()
	if err != nil {
		return nil, err
	}

	if !*res {
		return nil, errors.New("printer is not giving frame")
	}

	data := <-ch
	close(ch)

	return &data, nil
}

type rpcPutRawParams struct {
	FileID string `json:"file_id"`
	Length int    `json:"length"`
}

func (c *Client) sendFilePart(part *[]byte, id *string) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	err := c.call("put_raw", rpcPutRawParams{*id, len(*part)}, nil)
	if err != nil {
		return err
	}

	_, err = c.rpc.Write(*part)
	if err != nil {
		return err
	}

	return nil
}

type rpcPutInitParams struct {
	BlockSize int    `json:"block_size"`
	FileID    string `json:"file_id"`
	FilePath  string `json:"file_path"`
	Length    int    `json:"length"`
}

type rpcPutTermParams struct {
	Checksum uint32 `json:"crc"`
	FileID   string `json:"file_id"`
	Length   int    `json:"length"`
}

// PutFile sends a file to the printer and saves at the specified remote path
func (c *Client) PutFile(path string, r io.ReadCloser, size int) error {
	fileID := uuid.New().String()

	err := c.call("put_init", rpcPutInitParams{
		BlockSize: printFileBlockSize,
		FileID:    fileID,
		FilePath:  path,
		Length:    size,
	}, nil)
	if err != nil {
		return err
	}

	checksum := crc32.NewIEEE()

	bs := make([]byte, printFileBlockSize)

	for i := 0; i < size; i += printFileBlockSize {
		_, err := r.Read(bs)
		if err != nil {
			return err
		}

		checksum.Write(bs)

		c.sendFilePart(&bs, &fileID)
		if err != nil {
			return err
		}
	}

	bs = nil // explicitly deref

	return c.call("put_term", rpcPutTermParams{checksum.Sum32(), fileID, size}, nil)
}

type rpcPrintParams struct {
	FilePath     string `json:"filepath"`
	TransferWait bool   `json:"transfer_wait"`
}

// Print will synchronously print a .makerbot file with the provided
// `filename` (can be anything). `data` should be the contents of the
// .makerbot file. The function returns when it is done sending the entire
// file. If you want to monitor progress of the upload, see HandleStateChange.
//
// For easier usage, see PrintFile.
func (c *Client) Print(filename string, r io.ReadCloser, size int) error {
	err := c.call("print", rpcPrintParams{filename, true}, nil)
	if err != nil {
		return err
	}

	err = c.call("process_method", rpcProcessMethodParams{"build_plate_cleared"}, nil)
	if err != nil {
		return err
	}

	return c.PutFile(fmt.Sprintf("/current_thing/%s", filename), r, size)
}

// PrintFile is a convenience method for Print, taking in a
// `filename` and automatically reading from it then
// feeding it to Print.
func (c *Client) PrintFile(filename string) error {
	fil, err := os.Open(filename)
	if err != nil {
		return err
	}

	stat, err := os.Stat(filename)
	if err != nil {
		return err
	}

	return c.Print(filepath.Base(filename), fil, int(stat.Size()))
}

// PrintFileVerify is exactly like PrintFile except it errors
// if the print file is not designed for the printer that this
// Client is connected to.
func (c *Client) PrintFileVerify(filename string) error {
	metadata, err := printfile.GetFileMetadata(filename)
	if err != nil {
		return err
	}

	if metadata.BotType != c.Printer.BotType {
		return fmt.Errorf("print file was not sliced for this MakerBot printer (got: %s, wanted: %s)", metadata.BotType, c.Printer.BotType)
	}

	return c.PrintFile(filename)
}

type rpcCopySSHIDParams struct {
	FilePath string `json:"filepath"`
}

func (c *Client) copySSHID(path string) error {
	return c.call("copy_ssh_id", rpcCopySSHIDParams{path}, nil)
}

type rpcSetStagingURLsParams struct {
	ReflectorURL   string `json:"reflector_url"`
	ThingiverseURL string `json:"thingiverse_url"`
}

// SetStagingURLs points the bot to arbitrary URLs for its web services.
func (c *Client) SetStagingURLs(reflectorURL, thingiverseURL string) error {
	return c.call("set_staging_urls", rpcSetStagingURLsParams{reflectorURL, thingiverseURL}, nil)
}

type rpcAddMakerBotAccountParams struct {
	Username      string `json:"username"`
	MakerBotToken string `json:"makerbot_token"`
}

// AddMakerBotAccount authorizes a MakerBot account to the printer
func (c *Client) AddMakerBotAccount(username, token string) error {
	return c.call("add_makerbot_account", rpcAddMakerBotAccountParams{username, token}, nil)
}

// CopySSHPublicKey copies an SSH public key to the printer, allowing
// one to SSH into the printer as `root` with the key.
//
// EXPERIMENTAL: may not work!
func (c *Client) CopySSHPublicKey(key []byte) error {
	r := ioutil.NopCloser(bytes.NewReader(key))
	fp := fmt.Sprintf("%d", time.Now().UnixNano())

	err := c.PutFile(fp, r, len(key))
	if err != nil {
		return err
	}

	return c.copySSHID(fp)
}
