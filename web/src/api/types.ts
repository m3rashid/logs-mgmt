// Enums matching Go types
export type AircraftStatus = 0 | 1 | 2 // 0: Inactive, 1: Active, 2: Maintenance
export type FlightStatus = 0 | 1 | 2 | 3 // 0: Scheduled, 1: InFlight, 2: Completed, 3: Failed
export type LogLevel = 0 | 1 | 2 | 3 | 4 // 0: None, 1: Info, 2: Warning, 3: Error, 4: Fatal

// Models matching Go structs (with snake_case JSON fields)
export type Aircraft = {
  id: number
  name: string
  model: string
  serial_number: string
  battery_health: number
  status: AircraftStatus
  last_maintenance_date: string
  maintenance_due_date: string
  created_at: string
  updated_at: string
}

export type Flight = {
  id: number
  aircraft_id: number
  aircraft?: Aircraft
  source_lat: number
  source_lon: number
  destination_lat: number
  destination_lon: number
  start_time: string
  end_time: string
  current_location_lat: number
  current_location_lon: number
  status: FlightStatus
}

export type SystemInfo = {
  cpu: number
  memory: number
  disk: number
  network: number
  battery: number
  uptime: number
}

export type FlightEventData = {
  speed: number
  lat: number
  lon: number
  alt: number
}

export type FlightEvent = {
  id: string
  flight_id: number
  aircraft_id: number
  time: string
  level: LogLevel
  data: FlightEventData
  sys_info: SystemInfo
}

// Input types for create operations
export type CreateAircraftInput = {
  name: string
  model: string
  serial_number: string
  battery_health: number
  status: AircraftStatus
  last_maintenance_date?: string
  maintenance_due_date?: string
}

export type CreateFlightInput = {
  aircraft_id: number
  source_lat: number
  source_lon: number
  destination_lat: number
  destination_lon: number
  start_time: string
  end_time: string
}
