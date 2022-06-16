package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"pretty"
	"runtime"
	"strconv"
	"strings"
)

var inputInteractive io.Reader = os.Stdin
var outStream io.Writer = os.Stdout

func main() {

	//Input file or interactive mode
	ii := len(os.Args)
	var scanner *bufio.Scanner
	switch {
	case ii > 4:
		log.Fatal("Unknown command line input")
	case ii == 2:
		inputFile, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer inputFile.Close()
		scanner = bufio.NewScanner(inputFile)
	default:
		scanner = bufio.NewScanner(inputInteractive)
	}

	//Create a carpark
	var carpark = &Carpark{}

	//Operate the carpark
	operateCarpark(carpark, scanner)
}

//operateCarpark reads input queries from console or text file and executes the command
func operateCarpark(carpark *Carpark, scanner *bufio.Scanner) {
	newlineStr := getNewlineStr()
	exit := false
	for !exit && scanner.Scan() {
		input := scanner.Text()
		input = strings.TrimRight(input, newlineStr)
		s := parse(input)

		switch {
		case s[0] == "Create_parking_lot" && len(s) == 2: //Initialize carpark
			maxSlot, err := strconv.Atoi(s[1])
			if checkError(err) {
				break
			}
			err = carpark.init(maxSlot)
			if !checkError(err) {
				fmt.Fprintf(outStream, "Created parking of %v slots\n", maxSlot)
			}

		case s[0] == "Park" && len(s) == 4: //Park a new car
			age,_ := strconv.Atoi(s[3]);
      car := Car{
				registration: s[1],
				driver_age:       age,
			}
			slotNo, err := carpark.insertCar(&car)
			if !checkError(err) {
				fmt.Fprintf(outStream, "Car with vehicle registration number %s has been parked at slot number %v\n",car.registration, slotNo)
			}

		case s[0] == "Leave" && len(s) == 2: //Remove a parked car
			slotNo, err := strconv.Atoi(s[1])
			if checkError(err) {
				break
			}
 
			car, err := carpark.removeCar(slotNo)
			if !checkError(err) {
				fmt.Fprintf(outStream, "Slot number %v is vacated, the car with vehicle registration number %s left the space, the driver of the car was of age %v \n",slotNo,car.registration,car.driver_age )
			}

		case s[0] == "Vehicle_registration_number_for_driver_of_age" && len(s) == 2: //Return slot numbers with given driver age
      age,_ := strconv.Atoi(s[1])
			_, registration, err := carpark.getCarsWithAge(age)
			if checkError(err) {
				break
			}
			err = pretty.Printer(registration, outStream)
			if err != nil {
				panic(err.Error())
			}

      case s[0] == "Slot_numbers_for_driver_of_age" && len(s) == 2: //Return slot numbers with given driver age
      age,_ := strconv.Atoi(s[1])
			slots, _, err := carpark.getCarsWithAge(age)
			if checkError(err) {
				break
			}
			err = pretty.Printer(slots, outStream)
			if err != nil {
				panic(err.Error())
			}

      

		case s[0] == "Slot_number_for_car_with_number" && len(s) == 2: //Return slot numbers with given car registration number
			slotNo, err := carpark.getCarWithRegistrationNo(s[1])
			if !checkError(err) {
				fmt.Fprintln(outStream, slotNo)
			}

		default: //Default option
			fmt.Fprintln(outStream, "Unknown input command")
		}
	}
}

//getNewlineStr identifies operating system and returns newline character used
func getNewlineStr() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func parse(input string) []string {
	s := strings.Split(input, " ")
	return s
}

func checkError(err error) bool {
	if err != nil {
		fmt.Fprintln(outStream, err.Error())
		return true
	}
	return false
}