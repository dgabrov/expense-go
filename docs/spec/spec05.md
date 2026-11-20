@internal/impl/util.go

- create function called getDateRange
- params int startMonth int startYear int endMonth int endYear
- the values must be valid - otherwise return error
- returns two *Date and one error
- first date is inclusive (and it is start date of the start month)
- start month and year are mandatory
- end month and year are not mandatory, but either both are present and valid or none otherwise error
- start date is inclusive, end date is exclusive. So if I put end date 12 2025 you return for end date 2026-01-01 and I will see to it to seek dates smaller than this
