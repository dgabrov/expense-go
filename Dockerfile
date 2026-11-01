FROM amd64/alpine
EXPOSE 3001
WORKDIR /app
COPY expense_backend .
CMD ["/app/expense_backend"]
