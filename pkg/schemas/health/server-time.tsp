import "@typespec/http";
import "@typespec/rest";
import "./health.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using HealthAPI;

/**
 * Core server time information with timezone details
 */
model ServerTimeInfo {
  /** Server timezone identifier (e.g., "UTC", "America/New_York") */
  timezone: string;
  
  /** Current time in UTC */
  utcTime: utcDateTime;
  
  /** Current time in server's local timezone */
  localTime: utcDateTime;
  
  /** UTC offset string (e.g., "+00:00", "-05:00") */
  utcOffset: string;
  
  /** Unix timestamp in seconds */
  epoch: int64;
  
  /** ISO8601 formatted timestamp */
  iso8601: string;
  
  /** RFC3339 nano precision timestamp */
  rfc3339Nano: string;
  
  /** Timezone abbreviation (e.g., "UTC", "EST") */
  timezoneAbbr: string;
  
  /** Whether daylight saving time is active */
  isDST: boolean;
}

/**
 * Various formatted timestamp representations
 */
model FormattedTimestamps {
  /** Human-readable format (e.g., "Monday, January 2, 2006 at 3:04:05 PM MST") */
  human: string;
  
  /** Date only in YYYY-MM-DD format */
  date: string;
  
  /** Time only in HH:MM:SS format */
  time: string;
  
  /** DateTime in YYYY-MM-DD HH:MM:SS format */
  datetime: string;
  
  /** RFC822 formatted timestamp */
  rfc822: string;
  
  /** RFC850 formatted timestamp */
  rfc850: string;
  
  /** RFC1123 formatted timestamp */
  rfc1123: string;
  
  /** Kitchen format (3:04PM) */
  kitchen: string;
  
  /** Stamp format (Jan _2 15:04:05) */
  stamp: string;
}

/**
 * Unix timestamps in different precision levels
 */
model UnixTimestamps {
  /** Unix timestamp in seconds */
  seconds: int64;
  
  /** Unix timestamp in milliseconds */
  milliseconds: int64;
  
  /** Unix timestamp in microseconds */
  microseconds: int64;
  
  /** Unix timestamp in nanoseconds */
  nanoseconds: int64;
}

/**
 * Timezone information details
 */
model TimezoneInfo {
  /** Full timezone name */
  name: string;
  
  /** Timezone abbreviation */
  abbreviation: string;
  
  /** Offset from UTC in seconds */
  offsetSeconds: int32;
  
  /** Offset from UTC in hours (decimal) */
  offsetHours: float64;
  
  /** Whether daylight saving time is currently active */
  isDST: boolean;
  
  /** Standard timezone abbreviation */
  standardAbbr: string;
  
  /** Daylight saving timezone abbreviation */
  daylightAbbr?: string;
}

/**
 * Complete server time response with all timestamp formats
 */
model ServerTime {
  /** Core server time information */
  serverTime: ServerTimeInfo;
  
  /** Various formatted timestamp representations */
  formatted: FormattedTimestamps;
  
  /** Unix timestamps in different precisions */
  unixTimestamps: UnixTimestamps;
  
  /** Detailed timezone information */
  timezoneInfo: TimezoneInfo;
  
  /** Optional OpenTelemetry trace ID for correlation */
  traceId?: string;
  
  /** Optional Server-Timing metrics */
  serverTiming?: ServerTimingMetrics;
}
