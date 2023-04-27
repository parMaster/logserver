# Logserver 

This service can be connected to MQTT broker, subscribed to the topics where telemetry data is reported, store the data in a database and visualize data with charts.

No additional documentation will be provided, code is self-explanatory and trivial.

This was ment to be an educational project to learn Golang, so a couple words about what I learned here:

  - working with MQTT broker and subscribing to topics (using Paho MQTT library). Getting used to the idea of asynchronous programming
	
  - using standart structure for Golang projects - which is overkill for this project, but I wanted to learn it
	
  - using different data stores via interfaces
	
  - using SQLite as a database. In the first version I used PostgreSQL, but it was an overkill for this project (I was surprised finding out that you need a stored procedure to get the last inserted row ID or trigger to implemented autoincrement in PostgreSQL.. crazy stuff..), although it was a good learning experience, specially using it with Docker.

Currently it's working on a local Raspberry Pi 4, connected to a hosted MQTT broker and used mainly to monitor the climate in the crocodile enclosure. Climate control in the enclosure is done by a custom made ESP32-based device, which is also connected to the MQTT broker and reports the data - source is [available in my other repository](https://github.com/parMaster/ESP32Base).

Frontend contains graphs to display the data and in my case it looks like this:

![Logserver_Front](https://user-images.githubusercontent.com/1956191/234888928-e4ea6679-256f-49f0-ba45-8df32f15e90e.jpg)
