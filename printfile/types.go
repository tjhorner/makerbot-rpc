package printfile

// MakerBotFile is a representation of a .makerbot file
type MakerBotFile struct {
	ThumbnailSizes *[]Thumbnail
	Toolpath       *Toolpath
	Metadata       *Metadata
}

// Thumbnail represents a `thumbnail_*.jpg` file embedded
// within a .makerbot file.
//
// TargetWidth and TargetHeight exist because of a really funny
// bug with MakerBot Print: since it takes the thumbnail photos
// on the device that is slicing the print and they didn't take
// into account the screen resolution/density, sometimes (e.g. on
// a MacBook Pro) it can be different than the target width that's
// stated on the file itself. On my MacBook, it's actually double
// the size lol.
//
// So ActualWidth and ActualHeight hold the actual dimensions of the
// image.
type Thumbnail struct {
	Data         []byte
	TargetWidth  int
	TargetHeight int
	ActualWidth  int
	ActualHeight int
}

// Toolpath is a set of ToolpathInstructions
type Toolpath []ToolpathInstruction

// ToolpathInstruction represents something
type ToolpathInstruction struct {
	Command ToolpathCommand `json:"command"`
}

// ToolpathCommand is a command
type ToolpathCommand struct {
	Function string `json:"function"`
	Metadata struct {
		Relative struct {
			A bool `json:"a"`
			X bool `json:"x"`
			Y bool `json:"y"`
			Z bool `json:"z"`
		} `json:"relative"`
	} `json:"metadata"`
	Parameters struct {
		A        float64 `json:"a"`
		FeedRate float64 `json:"feedrate"`
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
		Z        float64 `json:"z"`
	} `json:"parameters"`
	Tags []string `json:"tags"`
}

// Metadata is a representation of the meta.json
// file inside of .makerbot print files.
//
// I have no clue what most of these fields do.
// I just threw meta.json into https://mholt.github.io/json-to-go/
// because I was NOT doing all of this by hand.
type Metadata struct {
	BotType     string `json:"bot_type"`
	BoundingBox struct {
		XMax float64 `json:"x_max"`
		XMin float64 `json:"x_min"`
		YMax float64 `json:"y_max"`
		YMin float64 `json:"y_min"`
		ZMax float64 `json:"z_max"`
		ZMin float64 `json:"z_min"`
	} `json:"bounding_box"`
	ChamberTemperature       float64   `json:"chamber_temperature"`
	CommandedDurationSeconds float64   `json:"commanded_duration_s"`
	DurationSeconds          float64   `json:"duration_s"`
	ExtruderTemperature      int       `json:"extruder_temperature"`
	ExtruderTemperatures     []int     `json:"extruder_temperatures"`
	ExtrusionDistanceMm      float64   `json:"extrusion_distance_mm"`
	ExtrusionDistancesMm     []float64 `json:"extrusion_distances_mm"`
	ExtrusionMassGrams       float64   `json:"extrusion_mass_g"`
	ExtrusionMassesGrams     []float64 `json:"extrusion_masses_g"`
	GrueVersion              string    `json:"grue_version"`
	MachineConfig            struct {
		Acceleration struct {
			BufferSize              int `json:"buffer_size"`
			ImpulseSpeedLimitMmPerS struct {
				X int `json:"x"`
				Y int `json:"y"`
				Z int `json:"z"`
			} `json:"impulse_speed_limit_mm_per_s"`
			MaxSpeedChangeMmPerS struct {
				X int `json:"x"`
				Y int `json:"y"`
				Z int `json:"z"`
			} `json:"max_speed_change_mm_per_s"`
			MinSpeedChangeMmPerS struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z int     `json:"z"`
			} `json:"min_speed_change_mm_per_s"`
			RateMmPerSSq struct {
				X int `json:"x"`
				Y int `json:"y"`
				Z int `json:"z"`
			} `json:"rate_mm_per_s_sq"`
			SplitMoveDistanceMm     float64 `json:"split_move_distance_mm"`
			SplitMoveRecursionCount int     `json:"split_move_recursion_count"`
		} `json:"acceleration"`
		BotType     string `json:"bot_type"`
		BuildVolume struct {
			X int `json:"x"`
			Y int `json:"y"`
			Z int `json:"z"`
		} `json:"build_volume"`
		ExtraSlicerSettings struct {
			PlateVariability float64 `json:"plate_variability"`
		} `json:"extra_slicer_settings"`
		ExtruderProfiles struct {
			AttachedExtruders []struct {
				Calibrated string `json:"calibrated"`
				ID         int    `json:"id"`
			} `json:"attached_extruders"`
			Mk12 struct {
				Materials struct {
					Pla struct {
						Acceleration struct {
							ImpulseSpeedLimitMmPerS struct {
								A int `json:"a"`
							} `json:"impulse_speed_limit_mm_per_s"`
							MaxSpeedChangeMmPerS struct {
								A float64 `json:"a"`
							} `json:"max_speed_change_mm_per_s"`
							MinSpeedChangeMmPerS struct {
								A float64 `json:"a"`
							} `json:"min_speed_change_mm_per_s"`
							RateMmPerSSq struct {
								A int `json:"a"`
							} `json:"rate_mm_per_s_sq"`
							SlipCompensationTable [][]int `json:"slip_compensation_table"`
						} `json:"acceleration"`
						FeedDiameter          float64 `json:"feed_diameter"`
						MaxFlowRate           float64 `json:"max_flow_rate"`
						OozeFeedstockDistance float64 `json:"ooze_feedstock_distance"`
						RestartRate           int     `json:"restart_rate"`
						RetractDistance       float64 `json:"retract_distance"`
						RetractRate           int     `json:"retract_rate"`
						Temperature           int     `json:"temperature"`
					} `json:"pla"`
				} `json:"materials"`
				MaxSpeedMmPerSecond struct {
					A float64 `json:"a"`
				} `json:"max_speed_mm_per_second"`
				NozzleDiameter float64 `json:"nozzle_diameter"`
				StepsPerMm     struct {
					A float64 `json:"a"`
				} `json:"steps_per_mm"`
			} `json:"mk12"`
			Mk13 struct {
				Materials struct {
					Pla struct {
						Acceleration struct {
							ImpulseSpeedLimitMmPerS struct {
								A float64 `json:"a"`
							} `json:"impulse_speed_limit_mm_per_s"`
							MaxSpeedChangeMmPerS struct {
								A float64 `json:"a"`
							} `json:"max_speed_change_mm_per_s"`
							MinSpeedChangeMmPerS struct {
								A float64 `json:"a"`
							} `json:"min_speed_change_mm_per_s"`
							RateMmPerSSq struct {
								A float64 `json:"a"`
							} `json:"rate_mm_per_s_sq"`
							SlipCompensationTable [][]int `json:"slip_compensation_table"`
						} `json:"acceleration"`
						FeedDiameter          float64 `json:"feed_diameter"`
						MaxFlowRate           float64 `json:"max_flow_rate"`
						OozeFeedstockDistance float64 `json:"ooze_feedstock_distance"`
						RestartRate           int     `json:"restart_rate"`
						RetractDistance       float64 `json:"retract_distance"`
						RetractRate           int     `json:"retract_rate"`
						Temperature           int     `json:"temperature"`
					} `json:"pla"`
				} `json:"materials"`
				MaxSpeedMmPerSecond struct {
					A float64 `json:"a"`
				} `json:"max_speed_mm_per_second"`
				NozzleDiameter float64 `json:"nozzle_diameter"`
				StepsPerMm     struct {
					A float64 `json:"a"`
				} `json:"steps_per_mm"`
			} `json:"mk13"`
			Mk13Impla struct {
				Materials struct {
					ImPla struct {
						Acceleration struct {
							ImpulseSpeedLimitMmPerS struct {
								A float64 `json:"a"`
							} `json:"impulse_speed_limit_mm_per_s"`
							MaxSpeedChangeMmPerS struct {
								A float64 `json:"a"`
							} `json:"max_speed_change_mm_per_s"`
							MinSpeedChangeMmPerS struct {
								A float64 `json:"a"`
							} `json:"min_speed_change_mm_per_s"`
							RateMmPerSSq struct {
								A float64 `json:"a"`
							} `json:"rate_mm_per_s_sq"`
							SlipCompensationTable [][]int `json:"slip_compensation_table"`
						} `json:"acceleration"`
						FeedDiameter          float64 `json:"feed_diameter"`
						MaxFlowRate           float64 `json:"max_flow_rate"`
						OozeFeedstockDistance float64 `json:"ooze_feedstock_distance"`
						RestartRate           int     `json:"restart_rate"`
						RetractDistance       float64 `json:"retract_distance"`
						RetractRate           int     `json:"retract_rate"`
						Temperature           int     `json:"temperature"`
					} `json:"im-pla"`
				} `json:"materials"`
				MaxSpeedMmPerSecond struct {
					A float64 `json:"a"`
				} `json:"max_speed_mm_per_second"`
				NozzleDiameter float64 `json:"nozzle_diameter"`
				StepsPerMm     struct {
					A float64 `json:"a"`
				} `json:"steps_per_mm"`
			} `json:"mk13_impla"`
			SupportedExtruders struct {
				Num0  interface{} `json:"0"`
				Num1  string      `json:"1"`
				Num2  string      `json:"2"`
				Num3  string      `json:"3"`
				Num4  string      `json:"4"`
				Num5  string      `json:"5"`
				Num6  string      `json:"6"`
				Num7  string      `json:"7"`
				Num8  string      `json:"8"`
				Num9  string      `json:"9"`
				Num10 string      `json:"10"`
				Num11 string      `json:"11"`
				Num12 string      `json:"12"`
				Num13 string      `json:"13"`
				Num14 string      `json:"14"`
				Num99 string      `json:"99"`
			} `json:"supported_extruders"`
		} `json:"extruder_profiles"`
		GantryConfiguration struct {
			MaxFillSpeed       int `json:"max_fill_speed"`
			MaxInnerShellSpeed int `json:"max_inner_shell_speed"`
			MaxOuterShellSpeed int `json:"max_outer_shell_speed"`
			TravelSpeedXy      int `json:"travel_speed_xy"`
			TravelSpeedZ       int `json:"travel_speed_z"`
		} `json:"gantry_configuration"`
		MakerbotGeneration  int `json:"makerbot_generation"`
		MaxSpeedMmPerSecond struct {
			X int `json:"x"`
			Y int `json:"y"`
			Z int `json:"z"`
		} `json:"max_speed_mm_per_second"`
		StartPosition struct {
			X int     `json:"x"`
			Y int     `json:"y"`
			Z float64 `json:"z"`
		} `json:"start_position"`
		StepsPerMm struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
			Z int     `json:"z"`
		} `json:"steps_per_mm"`
		Version string `json:"version"`
	} `json:"machine_config"`
	Material      string   `json:"material"`
	Materials     []string `json:"materials"`
	MiracleConfig struct {
		Bot       string   `json:"_bot"`
		Extruders []string `json:"_extruders"`
		Materials []string `json:"_materials"`
		DoRaft    bool     `json:"doRaft"`
		Gaggles   struct {
			Default struct {
				AdjacentFillLeakyConnections        bool    `json:"adjacentFillLeakyConnections"`
				AdjacentFillLeakyDistanceRatio      float64 `json:"adjacentFillLeakyDistanceRatio"`
				BacklashEpsilon                     float64 `json:"backlashEpsilon"`
				BacklashFeedback                    float64 `json:"backlashFeedback"`
				BacklashX                           float64 `json:"backlashX"`
				BacklashY                           float64 `json:"backlashY"`
				BaseInsetDistanceMultiplier         float64 `json:"baseInsetDistanceMultiplier"`
				BaseLayerHeight                     float64 `json:"baseLayerHeight"`
				BaseLayerWidth                      float64 `json:"baseLayerWidth"`
				BaseNumberOfShells                  int     `json:"baseNumberOfShells"`
				BedZOffset                          int     `json:"bedZOffset"`
				BridgeAnchorMinimumLength           float64 `json:"bridgeAnchorMinimumLength"`
				BridgeAnchorWidth                   float64 `json:"bridgeAnchorWidth"`
				BridgeMaximumLength                 float64 `json:"bridgeMaximumLength"`
				BrimsBaseWidth                      float64 `json:"brimsBaseWidth"`
				BrimsModelOffset                    float64 `json:"brimsModelOffset"`
				BrimsOverlapWidth                   float64 `json:"brimsOverlapWidth"`
				Coarseness                          float64 `json:"coarseness"`
				ComputeVolumeLike210                bool    `json:"computeVolumeLike2_1_0"`
				DefaultExtruder                     int     `json:"defaultExtruder"`
				DefaultSupportMaterial              int     `json:"defaultSupportMaterial"`
				Description                         string  `json:"description"`
				DoBacklashCompensation              bool    `json:"doBacklashCompensation"`
				DoBreakawaySupport                  bool    `json:"doBreakawaySupport"`
				DoBridging                          bool    `json:"doBridging"`
				DoBrims                             bool    `json:"doBrims"`
				DoExponentialDeceleration           bool    `json:"doExponentialDeceleration"`
				DoExternalSpurs                     bool    `json:"doExternalSpurs"`
				DoFanCommand                        bool    `json:"doFanCommand"`
				DoFanModulation                     bool    `json:"doFanModulation"`
				DoFixedLayerStart                   bool    `json:"doFixedLayerStart"`
				DoFixedShellStart                   bool    `json:"doFixedShellStart"`
				DoInternalSpurs                     bool    `json:"doInternalSpurs"`
				DoMinfill                           bool    `json:"doMinfill"`
				DoMixedRaft                         bool    `json:"doMixedRaft"`
				DoMixedSupport                      bool    `json:"doMixedSupport"`
				DoNewPathPlanning                   bool    `json:"doNewPathPlanning"`
				DoPaddedBase                        bool    `json:"doPaddedBase"`
				DoRaft                              bool    `json:"doRaft"`
				DoRateLimit                         bool    `json:"doRateLimit"`
				DoSplitLongMoves                    bool    `json:"doSplitLongMoves"`
				DoSupport                           bool    `json:"doSupport"`
				DoSupportUnderBridges               bool    `json:"doSupportUnderBridges"`
				ExponentialDecelerationMinSpeed     float64 `json:"exponentialDecelerationMinSpeed"`
				ExponentialDecelerationRatio        float64 `json:"exponentialDecelerationRatio"`
				ExponentialDecelerationSegmentCount int     `json:"exponentialDecelerationSegmentCount"`
				ExtruderProfiles                    []struct {
					DefaultTemperature int `json:"defaultTemperature"`
					ExtrusionProfiles  struct {
						Bridges struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate float64 `json:"feedrate"`
						} `json:"bridges"`
						Brims struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate float64 `json:"feedrate"`
						} `json:"brims"`
						FirstModelLayer struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate float64 `json:"feedrate"`
						} `json:"firstModelLayer"`
						FloorSurfaceFills struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate int     `json:"feedrate"`
						} `json:"floorSurfaceFills"`
						Infill struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate int     `json:"feedrate"`
						} `json:"infill"`
						Insets struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate int     `json:"feedrate"`
						} `json:"insets"`
						Outlines struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate int     `json:"feedrate"`
						} `json:"outlines"`
						Purge struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate int     `json:"feedrate"`
						} `json:"purge"`
						Raft struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate float64 `json:"feedrate"`
						} `json:"raft"`
						RaftBase struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate float64 `json:"feedrate"`
						} `json:"raftBase"`
						RoofSurfaceFills struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate int     `json:"feedrate"`
						} `json:"roofSurfaceFills"`
						SparseRoofSurfaceFills struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate int     `json:"feedrate"`
						} `json:"sparseRoofSurfaceFills"`
						Spurs struct {
							FanSpeed float64 `json:"fanSpeed"`
							Feedrate int     `json:"feedrate"`
						} `json:"spurs"`
					} `json:"extrusionProfiles"`
					ExtrusionVolumeMultiplier float64 `json:"extrusionVolumeMultiplier"`
					FeedDiameter              float64 `json:"feedDiameter"`
					IdleTemperature           int     `json:"idleTemperature"`
					NozzleDiameter            float64 `json:"nozzleDiameter"`
					OozeFeedstockDistance     float64 `json:"oozeFeedstockDistance"`
					PreOozeFeedstockDistance  float64 `json:"preOozeFeedstockDistance"`
					RestartExtraDistance      float64 `json:"restartExtraDistance"`
					RestartRate               int     `json:"restartRate"`
					RetractDistance           float64 `json:"retractDistance"`
					RetractRate               int     `json:"retractRate"`
					ToolchangeRestartDistance float64 `json:"toolchangeRestartDistance"`
					ToolchangeRestartRate     float64 `json:"toolchangeRestartRate"`
					ToolchangeRetractDistance float64 `json:"toolchangeRetractDistance"`
					ToolchangeRetractRate     float64 `json:"toolchangeRetractRate"`
				} `json:"extruderProfiles"`
				FanDefaultSpeed                  float64 `json:"fanDefaultSpeed"`
				FanLayer                         int     `json:"fanLayer"`
				FanModulationThreshold           float64 `json:"fanModulationThreshold"`
				FanModulationWindow              float64 `json:"fanModulationWindow"`
				FixedLayerStartX                 int     `json:"fixedLayerStartX"`
				FixedLayerStartY                 int     `json:"fixedLayerStartY"`
				FixedShellStartDirection         int     `json:"fixedShellStartDirection"`
				FloorSolidThickness              int     `json:"floorSolidThickness"`
				FloorSurfaceThickness            float64 `json:"floorSurfaceThickness"`
				FloorThickness                   float64 `json:"floorThickness"`
				HorizontalInset                  int     `json:"horizontalInset"`
				InfillDensity                    float64 `json:"infillDensity"`
				InfillShellSpacingMultiplier     float64 `json:"infillShellSpacingMultiplier"`
				InsetDistanceMultiplier          float64 `json:"insetDistanceMultiplier"`
				LayerHeight                      float64 `json:"layerHeight"`
				LeakyConnectionsAdjacentDistance float64 `json:"leakyConnectionsAdjacentDistance"`
				MaxConnectionLength              float64 `json:"maxConnectionLength"`
				MaxSparseFillThickness           float64 `json:"maxSparseFillThickness"`
				MaxSpurWidth                     float64 `json:"maxSpurWidth"`
				MinLayerDuration                 float64 `json:"minLayerDuration"`
				MinLayerHeight                   float64 `json:"minLayerHeight"`
				MinRaftBaseGap                   float64 `json:"minRaftBaseGap"`
				MinSpeedMultiplier               float64 `json:"minSpeedMultiplier"`
				MinSpurLength                    float64 `json:"minSpurLength"`
				MinSpurWidth                     float64 `json:"minSpurWidth"`
				MinThickInfillImprovement        float64 `json:"minThickInfillImprovement"`
				MinimumMoveDistance              float64 `json:"minimumMoveDistance"`
				ModelFillProfiles                struct {
					Bridge struct {
						Density             float64 `json:"density"`
						OrientationInterval int     `json:"orientationInterval"`
						OrientationOffset   int     `json:"orientationOffset"`
						OrientationRange    int     `json:"orientationRange"`
						Pattern             string  `json:"pattern"`
					} `json:"bridge"`
					FloorSurface struct {
						Density             float64 `json:"density"`
						OrientationInterval int     `json:"orientationInterval"`
						OrientationRange    int     `json:"orientationRange"`
						Pattern             string  `json:"pattern"`
					} `json:"floor_surface"`
					RoofSurface struct {
						Density             float64 `json:"density"`
						OrientationInterval int     `json:"orientationInterval"`
						OrientationRange    int     `json:"orientationRange"`
						Pattern             string  `json:"pattern"`
					} `json:"roof_surface"`
					Solid struct {
						Density             float64 `json:"density"`
						OrientationInterval int     `json:"orientationInterval"`
						OrientationRange    int     `json:"orientationRange"`
						Pattern             string  `json:"pattern"`
					} `json:"solid"`
					Sparse struct {
						Density                  float64 `json:"density"`
						DiamondFillTurnDistance  float64 `json:"diamondFillTurnDistance"`
						DiamondFillZigZagOverlap float64 `json:"diamondFillZigZagOverlap"`
						OrientationInterval      int     `json:"orientationInterval"`
						OrientationOffset        int     `json:"orientationOffset"`
						OrientationRange         int     `json:"orientationRange"`
						Pattern                  string  `json:"pattern"`
					} `json:"sparse"`
					SparseRoofSurface struct {
						Density                  float64 `json:"density"`
						DiamondFillTurnDistance  float64 `json:"diamondFillTurnDistance"`
						DiamondFillZigZagOverlap float64 `json:"diamondFillZigZagOverlap"`
						OrientationInterval      int     `json:"orientationInterval"`
						OrientationOffset        int     `json:"orientationOffset"`
						OrientationRange         int     `json:"orientationRange"`
						Pattern                  string  `json:"pattern"`
					} `json:"sparse_roof_surface"`
				} `json:"modelFillProfiles"`
				NumberOfBrims                        int           `json:"numberOfBrims"`
				NumberOfExtentShells                 int           `json:"numberOfExtentShells"`
				NumberOfInternalBrims                int           `json:"numberOfInternalBrims"`
				NumberOfShells                       int           `json:"numberOfShells"`
				NumberOfSparseShells                 int           `json:"numberOfSparseShells"`
				NumberOfSupportShells                int           `json:"numberOfSupportShells"`
				PaddedBaseOutlineOffset              float64       `json:"paddedBaseOutlineOffset"`
				PauseHeights                         []interface{} `json:"pauseHeights"`
				PurgeBaseRotation                    int           `json:"purgeBaseRotation"`
				PurgeBucketSide                      float64       `json:"purgeBucketSide"`
				PurgeWallBaseFilamentWidth           float64       `json:"purgeWallBaseFilamentWidth"`
				PurgeWallBasePatternLength           float64       `json:"purgeWallBasePatternLength"`
				PurgeWallBasePatternWidth            float64       `json:"purgeWallBasePatternWidth"`
				PurgeWallModelOffset                 float64       `json:"purgeWallModelOffset"`
				PurgeWallPatternWidth                float64       `json:"purgeWallPatternWidth"`
				PurgeWallSpacing                     float64       `json:"purgeWallSpacing"`
				PurgeWallWidth                       float64       `json:"purgeWallWidth"`
				PurgeWallXLength                     int           `json:"purgeWallXLength"`
				RaftBaseInfillShellSpacingMultiplier float64       `json:"raftBaseInfillShellSpacingMultiplier"`
				RaftBaseInsetDistanceMultiplier      float64       `json:"raftBaseInsetDistanceMultiplier"`
				RaftBaseLayers                       int           `json:"raftBaseLayers"`
				RaftBaseOutset                       int           `json:"raftBaseOutset"`
				RaftBaseShells                       int           `json:"raftBaseShells"`
				RaftBaseThickness                    float64       `json:"raftBaseThickness"`
				RaftBaseWidth                        float64       `json:"raftBaseWidth"`
				RaftBrimsSpacing                     float64       `json:"raftBrimsSpacing"`
				RaftExtraOffset                      float64       `json:"raftExtraOffset"`
				RaftFillProfiles                     struct {
					Base struct {
						Density                float64 `json:"density"`
						LinearFillGroupDensity float64 `json:"linearFillGroupDensity"`
						LinearFillGroupSize    int     `json:"linearFillGroupSize"`
						OrientationInterval    int     `json:"orientationInterval"`
						OrientationOffset      int     `json:"orientationOffset"`
						OrientationRange       int     `json:"orientationRange"`
						Pattern                string  `json:"pattern"`
					} `json:"base"`
					Interface struct {
						Density             float64 `json:"density"`
						OrientationInterval int     `json:"orientationInterval"`
						OrientationOffset   int     `json:"orientationOffset"`
						OrientationRange    int     `json:"orientationRange"`
						Pattern             string  `json:"pattern"`
					} `json:"interface"`
					Surface struct {
						Density             float64 `json:"density"`
						OrientationInterval int     `json:"orientationInterval"`
						OrientationOffset   int     `json:"orientationOffset"`
						OrientationRange    int     `json:"orientationRange"`
						Pattern             string  `json:"pattern"`
					} `json:"surface"`
				} `json:"raftFillProfiles"`
				RaftInterfaceLayers                int     `json:"raftInterfaceLayers"`
				RaftInterfaceShells                int     `json:"raftInterfaceShells"`
				RaftInterfaceThickness             float64 `json:"raftInterfaceThickness"`
				RaftInterfaceWidth                 float64 `json:"raftInterfaceWidth"`
				RaftInterfaceZOffset               float64 `json:"raftInterfaceZOffset"`
				RaftModelShellsSpacing             float64 `json:"raftModelShellsSpacing"`
				RaftModelSpacing                   float64 `json:"raftModelSpacing"`
				RaftSurfaceInsetDistanceMultiplier float64 `json:"raftSurfaceInsetDistanceMultiplier"`
				RaftSurfaceLayers                  int     `json:"raftSurfaceLayers"`
				RaftSurfaceOutset                  int     `json:"raftSurfaceOutset"`
				RaftSurfaceShellSpacingMultiplier  float64 `json:"raftSurfaceShellSpacingMultiplier"`
				RaftSurfaceShells                  int     `json:"raftSurfaceShells"`
				RaftSurfaceThickness               float64 `json:"raftSurfaceThickness"`
				RaftSurfaceZOffset                 float64 `json:"raftSurfaceZOffset"`
				RateLimitBufferSize                int     `json:"rateLimitBufferSize"`
				RateLimitMinSpeed                  int     `json:"rateLimitMinSpeed"`
				RateLimitSpeedRatio                float64 `json:"rateLimitSpeedRatio"`
				RateLimitSplitBias                 int     `json:"rateLimitSplitBias"`
				RateLimitSplitMoveDistance         float64 `json:"rateLimitSplitMoveDistance"`
				RateLimitSplitRecursionDepth       int     `json:"rateLimitSplitRecursionDepth"`
				RateLimitTransmissionRate          int     `json:"rateLimitTransmissionRate"`
				RoofAnchorMargin                   float64 `json:"roofAnchorMargin"`
				RoofSolidThickness                 int     `json:"roofSolidThickness"`
				RoofSurfaceThickness               float64 `json:"roofSurfaceThickness"`
				RoofThickness                      float64 `json:"roofThickness"`
				ShellsLeakyConnections             bool    `json:"shellsLeakyConnections"`
				SplitMinimumDistance               float64 `json:"splitMinimumDistance"`
				StartPosition                      struct {
					X int     `json:"x"`
					Y int     `json:"y"`
					Z float64 `json:"z"`
				} `json:"startPosition"`
				SupportAngle         float64 `json:"supportAngle"`
				SupportExtraDistance float64 `json:"supportExtraDistance"`
				SupportFillProfiles  struct {
					Sparse struct {
						ConsistentOrder     bool    `json:"consistentOrder"`
						Density             float64 `json:"density"`
						OrientationInterval int     `json:"orientationInterval"`
						OrientationOffset   int     `json:"orientationOffset"`
						OrientationRange    int     `json:"orientationRange"`
						Pattern             string  `json:"pattern"`
					} `json:"sparse"`
				} `json:"supportFillProfiles"`
				SupportInsetDistanceMultiplier float64 `json:"supportInsetDistanceMultiplier"`
				SupportInteriorExtruder        int     `json:"supportInteriorExtruder"`
				SupportLayerHeight             float64 `json:"supportLayerHeight"`
				SupportLeakyConnections        bool    `json:"supportLeakyConnections"`
				SupportModelSpacing            float64 `json:"supportModelSpacing"`
				SupportRoofModelSpacing        float64 `json:"supportRoofModelSpacing"`
				SupportShellSpacingMultiplier  float64 `json:"supportShellSpacingMultiplier"`
				ThickLayerThreshold            int     `json:"thickLayerThreshold"`
				ThickLayerVolumeMultiplier     int     `json:"thickLayerVolumeMultiplier"`
				TravelSpeedXY                  int     `json:"travelSpeedXY"`
				TravelSpeedZ                   int     `json:"travelSpeedZ"`
				UseRelativeExtruderPositions   bool    `json:"useRelativeExtruderPositions"`
			} `json:"default"`
		} `json:"gaggles"`
		Version string `json:"version"`
	} `json:"miracle_config"`
	ModelCounts []struct {
		Count int    `json:"count"`
		Name  string `json:"name"`
	} `json:"model_counts"`
	NumZLayers          int `json:"num_z_layers"`
	NumZTransitions     int `json:"num_z_transitions"`
	PlatformTemperature int `json:"platform_temperature"`
	Preferences         struct {
		Default struct {
			Overrides struct {
				DefaultSupportMaterial int `json:"defaultSupportMaterial"`
			} `json:"overrides"`
			PrintMode string `json:"print_mode"`
		} `json:"default"`
	} `json:"preferences"`
	ThingID       interface{} `json:"thing_id"`
	ToolType      string      `json:"tool_type"`
	ToolTypes     []string    `json:"tool_types"`
	TotalCommands int         `json:"total_commands"`
	UUID          string      `json:"uuid"`
	Version       string      `json:"version"`
}
