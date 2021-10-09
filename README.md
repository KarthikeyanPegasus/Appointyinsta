
# Appointy Assignment



## Instagram Backend API

Create an User
Should be a POST request 
Use JSON request body
URL should be ‘/users'

Get a user using id,
Should be a GET request,
Id should be in the url parameter,
URL should be ‘/users/<id here>’.

Create a Post,
Should be a POST request,
Use JSON request body,
URL should be ‘/posts'.

Get a post using id,
Should be a GET request,
Id should be in the url parameter,
URL should be ‘/posts/<id here>’.

List all posts of a user,
Should be a GET request,
URL should be ‘/posts/users/<Id here>'.

## Challenge Given 

We have to only use Standart Go library, MongoDb

"Challenge given for myself"

I have created both versions, with Mongodb, Without Mongodb or any other Databases.

## Advantages

Contains Unit Testing both versions, With MongoDB and Without MongoDB versions.

Password is secured using Message Digest 5.

The server is thread safe with the help mutex.

## Quality of codes Achieved
Code is reusable.

Consistency in naming.

All the endpoints are working in both versions



## FAQs
I)    While unit testing with Database version, Most of the testcases will return null / not found. Because the Database is not initialized in your system.

To rectify this error, You are requested to follow the instructions below.


<o>  Kindly replace the Database name in the code with the database that is available in your system.

<o>Kindly replace the collections name "users" with the collection which is having the schema of users, contains ID, Name, Email, Password.

<o> Kindly replace the collections name "Posts" with the collection which is having the schema of posts, contains Id, Caption, ImageURL, PTtime.

[Important. if a single word in the above FAQ changed, will result in fault working of the code](#)

ii)  curl returns null means, successfully connected to database, but either table is empty or table is not created successfully, retry the above step, make sure it is case-sensitive. 


## Submitted by

Karthikeyan P.

Flutter Developer [Link](https://drive.google.com/file/d/1ycUR-kfRzF5Jy390c6lvDwV0_Ym9_u57/view?usp=sharing)

Resume [Link](https://drive.google.com/file/d/1YdHaBioRBccwyQMDvZRQ1bgdNQKj4TIa/view?usp=sharing)
