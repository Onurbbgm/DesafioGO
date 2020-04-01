# GO Compare Data in CSV Files

This project was part of a challenge with the objective to compare two csv files and create a report with error messages pertinent to the problems found.

The program gets two csv files from a post send to it.

To build the project use the command ```go build``` in the folder with all the files.

To run use ```./DesafioGO```, this will start the server on the port 5000.

Then you can send the two files with a curl, like this one:

```
curl http://localhost:5000/ -F "fileOne=@clientData.csv" -F "fileTwo=@ourData.csv"
```

This will create a result file called resul.csv.
