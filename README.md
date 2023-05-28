Idea here is to build out some scripts in Golang as an alternative to the existing JS build.
Ultimately we want to test out concurrency to increase the efficiency of the webapp. 
At this point I'm building Mongo hooks for consistency with the existing DBM structure, this can be revised as well if needed.

Review comments on the main.go file on progress > at this time (5/28/23) the only thing published is the event button with an ID hook back to Mongo.
In my sandbox I'm testing out hashgen to see if we can increase call efficiency. 
As a next step I would like to interpret the user generation script into Golang, switch to hashgen and establish relationals between user and event IDs.
