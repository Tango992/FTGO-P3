# Ungraded Challenge 1

### Use a new database

```
use ungraded1DB
```

### Create into collections

```
db.books.insertMany([
    {_id: 1, title: "Kafka on the Shore", author: "Haruki Murakami", published: 2002, price: 25.00, stock: 35},
    {_id: 2, title: "Gadis Kretek", author: "Ratih Kumala", published: 2012, price: 14.25, stock: 100},
    {_id: 3, title: "Atomic Habits", author: "James Clear", published: 2018, price: 17.25, stock: 100},
    {_id: 4, title: "Norwegian Wood", author: "Haruki Murakami", published: 1987, price: 9.25, stock: 100},
])
```

### Read into collections

```
db.books.find()
```

### Update into collections

```
db.books.updateOne(
        {_id: 3},
        {$set: {stock: 250}}
)
```

### Delete from collections

```
db.books.deleteOne({_id: 4})
```

### Find book stock that are greater than 50

```
db.books.find({
    stock: {$gt: 50}
})
```

### Find book stock that are less than or equal to 40

```
db.books.find({
    stock: {$lte: 50}
})
```

### Find books that are published before 2010

```
db.books.find({
    published: {$lt: 2010}
})
```