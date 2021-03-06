<!DOCTYPE html>
<html>

<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <link rel="stylesheet" href="https://stackedit.io/style.css" />
</head>

<body class="stackedit">
  <div class="stackedit__html"><h1 id="todo-app">todo-app</h1>
<p><strong>TO-DO application in GO Language and MongoDB.</strong></p>
<p>The main purpose of this project is to get hands-on experience in Golang. There are very few tutorials and articles which gives you a complete end to end hands-on experience.<br>
In this project, we will build a todo app in which the server will be in Golang, the database will be MongoDB.</p>
<ul>
<li>Server — Go</li>
<li>Database — MongoDB</li>
</ul>
<p><strong>Assumption:</strong> Go is installed and have a basic understanding of it.</p>
<h2 id="lets-start"><strong>Let’s start</strong></h2>
<p>Create a project directory and give it an appropriate name.<br>
I am using <strong>todo-app</strong>.<br>
  Let’s first create the <code>go-server</code></p>
<p><strong>Golang Server</strong><br>
  The <code>go-server</code> directory structure will be:</p>
<pre><code>go-todo
  - go-server
    - middleware
      - middleware.go
    - models 
      - models.go
    - router 
      - router.go
    - main.go
</code></pre>
<p>For server, we need 2 dependencies: the first to connect with MongoDB and the second to create RESTAPIs.</p>
<ol>
<li>Connect with MongoDB
<ul>
<li>To install official mogo DB driver <a href="https://github.com/mongodb/mongo-go-driver">MongoDB Go Driver</a> run the below command in the terminal or command window<br>
<code>go get go.mongodb.org/mongo-driver</code></li>
</ul>
</li>
<li>Create RESTAPIs
<ul>
<li>To install the <code>gorilla/mux</code>  <a href="https://github.com/gorilla/mux">package</a> for the router, run the below command in the terminal or command window. <code>mux</code> is one of the most popular packages for the router in the Golang.<br>
<code>go get -u github.com/gorilla/mux</code></li>
</ul>
</li>
</ol>
<p><strong>Models</strong><br>
create a new folder <code>models</code> in the <code>go-server</code> directory and and create a new file<code>models.go</code> and paste the below code.</p>
<pre><code>package models
import  "go.mongodb.org/mongo-driver/bson/primitive"
</br>
type ToDoList struct {
   	&nbsp;&nbsp;&nbsp;&nbsp;ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  	&nbsp;&nbsp;&nbsp;&nbsp;Task string  `json:"task,omitempty"`
  	&nbsp;&nbsp;&nbsp;&nbsp;Status bool  `json:"status,omitempty"`
}
</br>
type Response struct {
&nbsp;&nbsp;&nbsp;&nbsp;Message string  `json:"message"`
}
</code></pre>
<p>In the  <code>ToDoList</code>  we have 3 fields:</p>
<ol>
<li>ID: This objectID will be generated by the MongoDB</li>
<li>Task: The test</li>
<li>Status:  <code>true</code>  or  <code>false</code></li>
</ol>
<p>The type of id in MongoDB is  <code>Object(id)</code>.</p>
<p>In the <code>Response</code> we have only one field <code>Message</code>, We are sending this response when we want to return any message.</p>
<p><strong>Middleware</strong><br>
Create a new folder <code>middleware</code> in the <code>go-server</code> directory and create a new file <code>middleware.go</code>. This file will contain the following functions,</p>
<ol>
<li>init</li>
<li>GetTask</li>
<li>GetAllTask</li>
<li>CreateTask</li>
<li>CompleteTask</li>
<li>UndoTask</li>
<li>DeleteTask</li>
<li>DeleteAllTask</li>
</ol>
<p>Let me explain the functionality. All the functions will be used in  <code>router.go</code>  which we will be writing in some time.</p>
<ul>
<li><code>init()</code><strong>:</strong> runs only once throughout the program life. In the init function the connection to the MongoDB will be established.</li>
<li><code>GetTask</code><strong>:</strong> First it set the header to tackle the “Cross-Origin Request” issue and then it will call the <code>getTask()</code>  function. It uses a<code>bson</code> package to get the data from the MongoDB.  <code>_bson.M_</code> _is used where M is an unordered.</li>
<li><code>GetAllTask</code><strong>:</strong> This is same as <code>GetTask</code> only the thing is it will return all the task present in the collection.</li>
<li><code>CreateTask</code><strong>:</strong> It first decodes the request body and store in  <code>models.ToDoList</code> type. It is imported from a  <code>models</code>  package. Then, it will call  <code>createTask</code> function and insert the task into the collection. Currently only single task you can create at a time.</li>
<li><code>CompleteTask</code><strong>:</strong> It is a <code>PUT</code> request where it will update the task’s status according to task ID. To get the  <code>params</code>  from the URL, we are using  <code>mux</code>  package.</li>
<li><code>UndoTask</code><strong>:</strong> This is same as <code>CompleteTask</code>, it only updates the task’s status to  <code>false</code></li>
<li><code>DeleteTask</code><strong>:</strong> It is a <code>DELETE</code> request. First, it’ll get the task id from the URL and then call <code>deleteTask</code> . It will retrieve the  <code>ObjectID</code>  of the task and then it will delete the task by its id from the collection.</li>
<li><code>DeleteAllTask</code><strong>:</strong> It deletes all the tasks from the collection.</li>
</ul>
<p>The middleware is complete.</p>
<p><strong>Router</strong><br>
Create a  <code>router</code>  folder in the  <code>go-server</code>  directory and then create a new file  <code>router.go</code>  in it. Paste the below code in the file.</p>
<pre><code>package router
import (
	"go-server/middleware"
	"github.com/gorilla/mux"
)
func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/task/{id}", middleware.GetTask).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/task", middleware.GetAllTask).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/task", middleware.CreateTask).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/task/delete/{id}", middleware.DeleteTask).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/task/delete/all/", middleware.DeleteAllTask).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/task/complete/{id}", middleware.CompleteTask).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/task/undo/{id}", middleware.UndoTask).Methods("PUT", "OPTIONS")
	return router
}
</code></pre>
<p><strong>main.go</strong><br>
Create a <code>main.go</code> file in the <code>go-server</code> directory. Paste the below code in it.</p>
<pre><code>package main
import (
	"fmt"
	"go-server/router"
	"log"
	"net/http"
)
func main() {
	router := router.Router()
	fmt.Println("Server is listening on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", router))
}
</code></pre>
<p>Import  <code>net/http</code>  package to serve the routes at  <code>8081</code>  port and  <code>./router</code>to import  <code>router</code>  package.</p>
<p>Create an instance of  <code>router</code>  package.</p>
<pre><code>router := router.Router()
</code></pre>
<p>Serve/host the application on the  <code>8081</code>  port.</p>
<pre><code>http.ListenAndServe(":8081", router)
</code></pre>
<p>Open the terminal from the  <code>go-server</code>  directory and run the below command to serve the server.</p>
<pre><code>go run main.go
</code></pre>
<p>You’ll see the output below,</p>
<pre><code>C:\Users\USER\todo-app\go-server&gt;go run main.go
Connected to MongoDB...
Collection instance created...
Server is listening on port 8081...
</code></pre>
Inspired by https://medium.com/@schadokar blogs.
  
<h2 id="since-this-file-is-getting-bit-lengthy--i-have-pushed-the-mongodb-setup-file-and-sample-request-response-file-in-the-repository-itself."><strong>Since this file is getting bit lengthy, I have pushed the MongoDB Setup file and Sample request-response file in the repository itself.</strong></h2>
</div>
</body>

</html>
