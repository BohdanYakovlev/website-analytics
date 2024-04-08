# [Website Analytics](https://docs.google.com/document/d/1fpnCkMOV-CqydPIo0EnSrQ6pln4pniSxoGa_9uAgLQ8/edit#heading=h.jfgry0741f91) Console App

This console application is designed to analyze website analytics data provided in CSV files representing user visits to product pages on different days. The task involves identifying users who visited pages on both days and, specifically, those who visited a page on the second day that they hadn't visited on the first day.

## How to Run the Code

To run this console app, follow these steps:

1. Clone this repository to your local machine:

```
git clone <repository_url>
```

2. Navigate to the directory containing the Go code:

```
cd website-analytics
```

3. Ensure you have Go installed on your system. If not, you can download it from [here](https://golang.org/dl/).

4. Execute the Go code by running the following command:

```
go run main.go
```

5. Follow the on-screen instructions to provide the paths to your CSV files containing website analytics data for the first and second days.

6. Once the program finishes executing, it will display the users who visited some pages on both days and users who visited a page on the second day that they hadn't visited on the first day.

## Algorithms

### 1. RAM Method O(n):
   - Both day visitors:
      - Read both CSV files into memory.
      - During the reading search for visitors of both day, add them to ignore list and delete visitors from memory.
      - Return the ignore list.
   - New pages visited on second day:
      - Read the first day file into memory.
      - Read the second day file and identify pages that were not visited on the first day.
      - Return the results.
#### Efficiency:
   - Ignore list and removing unnecessary ones from memory for first task.
   - Data from the second file is not written to memory for second task.
   - Map ensures uniqueness among users and products.
#### Running time: 1-2 ms

### 2. Disk Method O(n):
   - Read both CSV files and create .txt file with name "userID-productID" into the corresponding folders.
   - For day folders search files with the same names and files which exist only in second day folder.
   - Return the results.
#### Efficiency:
   - Creating files does not overload RAM.
#### Running time: 60-90 ms

## Specifications of the executing machine:
   - CPU: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
   - RAM: 16Gb 33MHz
   - Disk: SSD 660P Series
   - OS: Ubuntu 20.04.6 LTS 64-bit

## Conclusion

This console application efficiently analyzes website analytics data to identify users who visited pages on both days and those who visited a page on the second day that they hadn't visited on the first day. RAM method loads all data into memory and processes it. It's suitable for datasets that can fit into RAM and offers faster processing. Disk method reads and writes data to disk, making it suitable for larger datasets that may not fit into memory. It might be slower due to disk I/O operations.

Feel free to reach out if you have any questions or encounter any issues while running the code.
