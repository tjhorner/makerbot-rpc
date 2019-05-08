package makerbot

import (
	"bytes"
	"encoding/json"
)

// Printer represents a connected printer
type Printer struct {
	MachineType        string `json:"machine_type"`         // The codename for this machine type
	Vid                int    `json:"vid"`                  // Vendor ID of the printer
	IP                 string `json:"ip"`                   // The local IP of this printer
	Pid                int    `json:"pid"`                  // Product ID of the printer
	APIVersion         string `json:"api_version"`          // API verison
	Serial             string `json:"iserial"`              // Serial number of the printer
	SSLPort            string `json:"ssl_port"`             // Port at which the HTTPS server can be accessed
	MachineName        string `json:"machine_name"`         // User-defined printer name
	MotorDriverVersion string `json:"motor_driver_version"` // Version number of the motor driver
	BotType            string `json:"bot_type"`             // Codename for the bot type
	Port               string `json:"port"`                 // JSON-RPC port (usually 9999)
	FirmwareVersion    struct {
		Major  int `json:"major"`
		Minor  int `json:"minor"`
		Bugfix int `json:"bugfix"`
		Build  int `json:"build"`
	} `json:"firmware_version"`
	Metadata *PrinterMetadata `json:"-"`
}

// PrinterMetadata is an object that the printer will send periodically
// and contains information about the printer itself
type PrinterMetadata struct {
	AutoUnload         string                `json:"auto_unload"`
	DisabledErrors     []interface{}         `json:"disabled_errors"`
	BotType            string                `json:"bot_type"`
	Sound              bool                  `json:"sound"`
	MachineName        string                `json:"machine_name"`
	CurrentProcess     *PrinterProcess       `json:"current_process"`
	APIVersion         string                `json:"api_version"`
	HasBeenConnectedTo bool                  `json:"has_been_connected_to"`
	IP                 string                `json:"ip"`
	Toolheads          map[string][]Toolhead `json:"toolheads"`
	MachineType        string                `json:"machine_type"`
	FirmwareVersion    struct {
		Major  int
		Minor  int
		Bugfix int
		Build  int
	} `json:"firmware_version"`
}

// PrinterProcess represents what the printer's current task is
type PrinterProcess struct {
	ID                  int                `json:"id"`
	Filename            *string            `json:"filename,omitempty"`
	Complete            bool               `json:"complete"`
	FilamentExtruded    float32            `json:"filament_extruded"`
	PrintTemperatures   map[string]float32 `json:"print_temperatures"`
	Name                string             `json:"name"`
	Filepath            *string            `json:"filepath"`
	Methods             []string           `json:"methods"`
	Username            *string            `json:"username,omitempty"`
	CanPrintAgain       *bool              `json:"can_print_again,omitempty"`
	Progress            *int               `json:"progress,omitempty"`
	Cancellable         bool               `json:"cancellable"`
	Step                PrintProcessStep   `json:"step"`
	StartTime           *epochTime         `json:"start_time"`
	ElapsedTime         epochTime          `json:"elapsed_time"`
	Cancelled           bool               `json:"cancelled"`
	ThingID             *int               `json:"thing_id,omitempty"`
	Reason              *string            `json:"reason,omitempty"`
	ToolIndex           *int               `json:"tool_index,omitempty"`
	TemperatureSettings *[]int             `json:"temperature_settings,omitempty"`
	// TODO: error field
}

// PrintProcessStep is an enum that represents a step that a PrinterProcess can go through
type PrintProcessStep int

func (s PrintProcessStep) String() string { return printProcessStepToString[s] }

// Humanize will return a human-readable string that describes the current step
func (s PrintProcessStep) Humanize() string { return printProcessStepToHumanString[s] }

const (
	// StepUnknown represents a step unknown to us that the printer returned
	StepUnknown PrintProcessStep = iota
	// StepInitializing means the printer is initializing the current print job (getting the toolheads ready, etc.)
	StepInitializing
	// StepInitialHeating means the printer is heating up the extruder(s) and getting it ready for final heating
	StepInitialHeating
	// StepFinalHeating means the printer is heating up the extruders(s) to their target temperature and will begin printing when done
	StepFinalHeating
	// StepCooling means the printer is cooling the extruder(s) (because it overheated, the print job was cancelled, etc.)
	StepCooling
	// StepHoming means the printer is finding the start position and is doing a small calibration step before printing
	StepHoming
	// StepPositionFound means the printer has found the start position and will begin printing very shortly
	StepPositionFound
	// StepPreheatingResuming means the printer was instructed to preheat the extruder(s) and is about to do so
	StepPreheatingResuming
	// StepCalibrating means the printer was instructed to begin calibration and is currently in the process of doing so
	StepCalibrating
	// StepPrinting means the printer is currently printing the current job
	StepPrinting
	// StepEndSequence means the printer is about to complete the current job
	StepEndSequence
	// StepCancelling means the user has requested that the current job be cancelled (or something went horribly wrong) and the printer is stopping everything
	StepCancelling
	// StepSuspending means the user has requested that the current job be paused ("suspended")
	StepSuspending
	// StepSuspended means the current job is paused ("suspended")
	StepSuspended
	// StepUnsuspending means the user has requested that the current job be unpaused ("unsuspended")
	StepUnsuspending
	// StepPreheatingLoading means the printer is preheating the extruder in preparation for filament loading
	StepPreheatingLoading
	// StepPreheatingUnloading means the printer is preheating the extruder in preparation for filament unloading
	StepPreheatingUnloading
	// StepLoadingFilament means the printer is in the middle of loading filament into the extruder
	StepLoadingFilament
	// StepUnloadingFilament means the printer is in the middle of unloading filament from the extruder
	StepUnloadingFilament
	// StepStoppingFilament means... the filament is being stopped, I guess?
	StepStoppingFilament
	// StepCleaningUp means the print was just completed and the printer is "cleaning up" (e.g. resetting the extruder position, moving the build plate back to normal, etc.)
	StepCleaningUp
	// StepClearBuildPlate means the build plate needs to be cleared ???
	StepClearBuildPlate
	// StepError means something went horribly wrong
	StepError
	// StepLoadingPrintTool means the extruder is currently being attached to the printer
	StepLoadingPrintTool
	// StepWaitingForFile means the printer is waiting for the file transfer from the remote host (laptop, phone, whatever) to begin
	StepWaitingForFile
	// StepTransfer means the printer is currently receiving the print file from the remote host
	StepTransfer
	// StepFailed means the print failed for some reason. You're probably out of filament, idiot
	StepFailed
	// StepCompleted means the print has completed. Yeet
	StepCompleted
	// StepHandlingRecoverableFilamentJam is an oddly specific step that means the printer is attempting to recover from a filament slip or jam without human intervention
	StepHandlingRecoverableFilamentJam
	// StepRunning means the thing is running
	StepRunning
)

var printProcessStepToHumanString = map[PrintProcessStep]string{
	StepUnknown:                        "Unknown",
	StepInitializing:                   "Initializing",
	StepInitialHeating:                 "Initial Heating",
	StepFinalHeating:                   "Final Heating",
	StepCooling:                        "Cooling",
	StepHoming:                         "Finding Position",
	StepPositionFound:                  "Position Found",
	StepPreheatingResuming:             "Resuming Pre-Heating",
	StepCalibrating:                    "Calibrating",
	StepPrinting:                       "Printing",
	StepEndSequence:                    "Cleaning Up",
	StepCancelling:                     "Cancelling",
	StepSuspending:                     "Suspending",
	StepSuspended:                      "Suspended",
	StepUnsuspending:                   "Unsuspending",
	StepPreheatingLoading:              "Preparing For Filament Loading",
	StepPreheatingUnloading:            "Preparing For Filament Unloading",
	StepLoadingFilament:                "Loading Filament",
	StepUnloadingFilament:              "Unloading Filament",
	StepStoppingFilament:               "Stopping Filament Loading/Unloading",
	StepCleaningUp:                     "Cleaning up",
	StepClearBuildPlate:                "Waiting For Clear Build Plate",
	StepError:                          "Error",
	StepLoadingPrintTool:               "Loading Print Tool",
	StepWaitingForFile:                 "Waiting For File",
	StepTransfer:                       "Transferring File",
	StepFailed:                         "Failed",
	StepCompleted:                      "Completed",
	StepHandlingRecoverableFilamentJam: "Attempting Filament Jam Recovery",
	StepRunning:                        "Running",
}

var printProcessStepToString = map[PrintProcessStep]string{
	StepUnknown:                        "", // just in case
	StepInitializing:                   "initializing",
	StepInitialHeating:                 "initial_heating",
	StepFinalHeating:                   "final_heating",
	StepCooling:                        "cooling",
	StepHoming:                         "homing",
	StepPositionFound:                  "position_found",
	StepPreheatingResuming:             "preheating_resuming",
	StepCalibrating:                    "calibrating",
	StepPrinting:                       "printing",
	StepEndSequence:                    "end_sequence",
	StepCancelling:                     "cancelling",
	StepSuspending:                     "suspending",
	StepSuspended:                      "suspended",
	StepUnsuspending:                   "unsuspending",
	StepPreheatingLoading:              "preheating_loading",
	StepPreheatingUnloading:            "preheating_unloading",
	StepLoadingFilament:                "loading_filament",
	StepUnloadingFilament:              "unloading_filament",
	StepStoppingFilament:               "stopping_filament",
	StepCleaningUp:                     "cleaning_up",
	StepClearBuildPlate:                "clear_build_plate",
	StepError:                          "error_step",
	StepLoadingPrintTool:               "loading_print_tool",
	StepWaitingForFile:                 "waiting_for_file",
	StepTransfer:                       "transfer",
	StepFailed:                         "failed",
	StepCompleted:                      "completed",
	StepHandlingRecoverableFilamentJam: "handling_recoverable_filament_jam",
	StepRunning:                        "running",
}

var printProcessStepToID = map[string]PrintProcessStep{
	"initializing":                      StepInitializing,
	"initial_heating":                   StepInitialHeating,
	"final_heating":                     StepFinalHeating,
	"cooling":                           StepCooling,
	"homing":                            StepHoming,
	"position_found":                    StepPositionFound,
	"preheating_resuming":               StepPreheatingResuming,
	"calibrating":                       StepCalibrating,
	"printing":                          StepPrinting,
	"end_sequence":                      StepEndSequence,
	"cancelling":                        StepCancelling,
	"suspending":                        StepSuspending,
	"suspended":                         StepSuspended,
	"unsuspending":                      StepUnsuspending,
	"preheating_loading":                StepPreheatingLoading,
	"preheating_unloading":              StepPreheatingUnloading,
	"loading_filament":                  StepLoadingFilament,
	"unloading_filament":                StepUnloadingFilament,
	"stopping_filament":                 StepStoppingFilament,
	"cleaning_up":                       StepCleaningUp,
	"clear_build_plate":                 StepClearBuildPlate,
	"error_step":                        StepError,
	"loading_print_tool":                StepLoadingPrintTool,
	"waiting_for_file":                  StepWaitingForFile,
	"transfer":                          StepTransfer,
	"failed":                            StepFailed,
	"completed":                         StepCompleted,
	"handling_recoverable_filament_jam": StepHandlingRecoverableFilamentJam,
	"running":                           StepRunning,
}

// MarshalJSON implements a JSON marshaler
func (s PrintProcessStep) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(printProcessStepToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements a JSON unmarshaler
func (s *PrintProcessStep) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*s = printProcessStepToID[j]
	return nil
}

// Toolhead represents a toolhead connected to a printer
// (e.g. a Smart Extruder+)
type Toolhead struct {
	Preheating         bool    `json:"preheating"`
	FilamentPresence   bool    `json:"filament_presence"`
	TargetTemperature  float32 `json:"target_temperature"`
	Error              int     `json:"error"`
	Index              int     `json:"index"`
	ToolPresent        bool    `json:"tool_present"`
	ToolID             int     `json:"tool_id"`
	CurrentTemperature float32 `json:"current_temperature"`
}
