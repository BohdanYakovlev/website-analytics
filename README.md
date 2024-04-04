# Website Analytics Console App

This console application is designed to analyze website analytics data provided in CSV files representing user visits to product pages on different days. The task involves identifying users who visited pages on both days and, specifically, those who visited a page on the second day that they hadn't visited on the first day.

## How to Run the Code

To run this console app, follow these steps:

1. Clone this repository to your local machine:

```
git clone <repository_url>
```

2. Navigate to the directory containing the Go code:

```
cd website-analytics-console-app
```

3. Ensure you have Go installed on your system. If not, you can download it from [here](https://golang.org/dl/).

4. Execute the Go code by running the following command:

```
go run main.go
```

5. Follow the on-screen instructions to provide the paths to your CSV files containing website analytics data for the first and second days.

6. Once the program finishes executing, it will display the users who visited some pages on both days and users who visited a page on the second day that they hadn't visited on the first day.

## Algorithm Explanation

### 1. **RAM Method**:
   - Read both CSV files into memory.
   - Process the data in memory to identify users who visited pages on both days and those who visited a page on the second day that they hadn't visited on the first day.
   - Return the results.

### 2. **Disk Method**:
   - Read CSV files sequentially from disk.
   - For each user and product visit:
     - Create separate files for users and products visited on each day.
     - Compare files to identify users who visited pages on both days and those who visited a page on the second day that they hadn't visited on the first day.
   - Return the results.

## Conclusion

This console application efficiently analyzes website analytics data to identify users who visited pages on both days and those who visited a page on the second day that they hadn't visited on the first day. It offers flexibility in handling datasets of various sizes and provides insights into user behavior on the website.

Feel free to reach out if you have any questions or encounter any issues while running the code.
You can now use this updated README.md file with the corrected code blocks.
