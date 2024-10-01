# Event Collector and CSV Exporter

This Go application collects events via HTTP and exports them to CSV files. It's designed to handle various event types and dynamically update CSV files based on the received data.

## Features

- Receives JSON events via HTTP POST requests
- Stores events in memory, grouped by event type
- Exports events to CSV files on demand
- Dynamically updates existing CSV files with new data
- Handles multiple event types, each with its own CSV file
- Command-line interface for saving data and quitting the application

## Requirements

- Go 1.15 or higher

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/yourusername/event-collector.git
   ```
2. Navigate to the project directory:
   ```
   cd event-collector
   ```

## Usage

1. Start the application:
   ```
   go run main.go
   ```

2. The server will start on port 3080. You should see the message:
   ```
   Server started on port 3080
   ```

3. Send events to the application using HTTP POST requests to `http://localhost:3080/events`. The payload should be a JSON array of event objects. Each event must have an `eventType` field. For example:
   ```json
   [
     {
       "eventType": "click",
       "elementId": "button1",
       "timestamp": 1623456789
     },
     {
       "eventType": "pageview",
       "url": "/home",
       "timestamp": 1623456790
     }
   ]
   ```

4. To save the collected events to CSV files, type `save-csv` in the terminal where the application is running. This will create or update CSV files in the `csv_exports` directory, one for each event type.

5. To quit the application, type `quit` in the terminal.

## CSV File Structure

- CSV files are named after the event type (e.g., `click.csv`, `pageview.csv`).
- Each CSV file includes all fields from all events of that type.
- If new fields are introduced in later events, they are added as new columns.
- Existing CSV files are updated with new data, preserving previous entries.

## Notes

- This application stores events in memory. For production use with large volumes of data, consider implementing a database backend.
- The application does not implement authentication or authorization. Ensure to add security measures for production use.
- CSV files are regenerated entirely on each save. For very large datasets, consider implementing an append-only approach.
