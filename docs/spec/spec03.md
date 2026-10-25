There is a project - REACTJS at location C:\Users\dgabr\work\expense-web - typescript
There is an interface here which is the facade for the backend
This project we work on, the current project, this one is the backend for that web application. 


Look at @C:\Users\dgabr\work\expense-web\src\server\servr.ts and basically except login which is already done, 
create endpoint for each of the operations that you see. ALL OF THEM POST.

Structure is similar with login: a new interface, a process function that returns payload or error and then use of writeJsonResponse. 

Inside the process function, you parse the request to the adequate type - types should be already in data package but if not ASK me first and then create the type

If you see interface type that extends another type, please do not embed the parent type in struct, rather copy the fields, there is only one situation about that

