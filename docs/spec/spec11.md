create new endpoint 

PUT /setCategory

payload of json type

{   
    transactionId: string
    categoryId: string 
}

basically for the transaction with the given Id sets the categoryId. Of course, referential integrity throws error. Also 0 records updated (meaning the transaction does not exist) also throws error. 

In case of success does not return anything but http 200
