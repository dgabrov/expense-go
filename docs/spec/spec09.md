create new endpoint 

GET /month/group

Query parameters
month - 1-12 mandatory
year - a valid year
groupId - the group id; can be empty string!!! in which case you get the ones without categoryId filled out

has security like the other ones with security (look at POST /addBulk)

returns slice of transactions
