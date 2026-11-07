implement @internal/proj/controller/post_report_last.go (currently has some hello world) as follows:

- get a request of type ReportRequest
- ignore ExcludeNone
- ignore End value 
- check Start is properly filled using the validate

Then prepare start date from Start of type data.Date, and run a query on the transaction table and get all the accounts in transactions that come after this date (inclusive) and the maximum transaction date.

create whatever method you need in server.Server

Return a response of type slice of type ReportLastResponse

