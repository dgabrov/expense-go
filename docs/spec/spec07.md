implement @internal/proj/controller/post_report_category.go (currently has some hello world) as follows:

- get a request of type ReportRequest
- Start is mandatory and must be validated
- ignore End value if null, meaning End attributes will not be part of any database queries
- anyway use getDateRange to get start and end dates
- ExcludeNone - exclude group that has the name "none". Effectively. The name not the id

- response is of type ReportCategoryResponse
- you select with adequate sql groups from transaction, filtering by transaction_date, grouping by categoryId 
- if CategoryId is empty or null please return its value, do not skip it into the json so that I see it, those are the not assigned transactions, actually coalesce to empty string works perfect, like you did for the change

- create whatever service.Service entries you need
- create http file with tests like you did for group report
