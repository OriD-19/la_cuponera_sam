# General API description
Make a separate API call for each category.

Each component will make a different API call to an endpoint which filters the coupons by catgory. This will result in more API calls, but in better performance overall, because we only need to specify a LIMIT query parameter inside the URL to receive the top 3 coupons. Or, if later on we want to change the amount of coupons, we can do so by modifying said LIMIT number.

Generate a UUID inside the frontend, and query the value of such ID when creating the element in the route.

For example, for creating  new element, just issue a PUT request at
`PUT "api/v1/coupons/{couponID}"`

# Data model
Store user attributes with the composite sort key model.

So, partition the data into the two main entities: **coupons and users.**

Users will be of different types, according to the project description.

Compose the sort key with the description of the value (for example, COUPON) followed by the ID
of the element we are searching for.
For example, we could use something like `COUPON#INFORMATION#1234` to retrieve the information
of the coupon with the ID 1234.
Then, we could just use a USERID partition key for each user to retrieve information about each coupon
for each user.

## Calling the Query API
When querying data, we could do something like this:

```
--key-condition-expression: id = :id
--expression-attribute-values '{
    ":id": {"S": "USER#${my_user_id}"}
}'
```

Also, for querying a particular sub-category of items, we could use the `begins_with()` expression 
function to check if an attribute begins with a given substring.

As you can see, we could just concatenate the user ID to the partition key information. In this case, since we want to retrive a USER, we could prepend the "USER" value to it and then append the specific USER ID that we want to check upon.

# Ordering the User idea

One of the most crucial decisions will be to choose how to store user details. 
Right now, we could simply use the same user entity for all types of user, and check for 
a specific attribute according to the type of user that we are looking for.

For example, if we want to log in as a client, we could simply put the client information
as a special type of attribute that is **only present inside client-type users**.
The same would apply for _enterprise-type_ users, or **administrator users**.