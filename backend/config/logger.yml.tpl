level: "{{ call . "LOG_LEVEL" }}"
encoding: "json"
encoderConfig:
  timeKey: "timestamp"
  levelKey: "severity"
  messageKey: "message"
  callerKey: "caller"
  levelEncoder: "capital"
  timeEncoder: "iso8601"
  durationEncoder: "string"
  callerEncoder: "short"
disableCaller: true
outputPaths:
  - "stderr"
errorOutputPaths:
  - "stderr"
