// Wrapper for functions exported by the lpo package.
// 01 - Jul. 12, 2018   First version, uploaded to github
// 02 - Sept. 6, 2018   Revised as first "delux" lporun now replacing old runopt


package main

import (
	"fmt"
	"github.com/go-opt/lpo"
	"github.com/pkg/errors"
	"strconv"
)


//==============================================================================

// wpGetPoint obtains a point (set of variable values for a specific constraint)
// that are used in other functions such as CalcLhs and CalcConViolation) and
// passes the row index and the list of values back to the caller via the argument
// list. In case of failure, function returns an error.
func wpGetPoint(rowIndex *int, point *[]float64) error {
	var userString string   // user input string
	var bigInt     int64    // big int translated from the user input
	var pointItem  float64  // item of the point list as the point is created
	var iElem      int      // index of element being processed
	var iCol       int      // index of column being processed
	var err        error    // error returned from functions called

	// Initialize the variables to be passed back, and prompt user for the value
	// of each variable.	
	*rowIndex = -1
	*point    = nil
	
	fmt.Printf("Enter constraint index: ")
	fmt.Scanln(&userString)
	if bigInt, err = strconv.ParseInt(userString,10,64); err != nil {
    	return errors.Errorf("'%s' is not an integer.", userString)
	}

	// Check that the input is valid, or fail with error if not. If valid, add
	// the variable value to the list.	
	*rowIndex = int(bigInt)

	if *rowIndex < 0 || *rowIndex >= len(lpo.Rows) {
		return errors.Errorf("Row index %d out of range.", *rowIndex)
	}

	fmt.Printf("Variable values are needed for the following constraint:\n")
	lpo.PrintRow(*rowIndex)	
	
	for i := 0; i < len(lpo.Rows[*rowIndex].HasElems); i++ {
		iElem      = lpo.Rows[*rowIndex].HasElems[i]
		iCol       = lpo.Elems[iElem].InCol
		userString = ""
		fmt.Printf("Enter value for %s: ", lpo.Cols[iCol].Name)
		fmt.Scanln(&userString)
		if pointItem, err = strconv.ParseFloat(userString, 64); err != nil {
    		return errors.Errorf("'%s' is not a real number.", userString)			
		}
		*point = append(*point, pointItem)	
	}
	
	return nil	
}

//==============================================================================

// wpCalcLhs is a wrapper for lpo.CalcLhs. It uses a separate function to get the
// constraint index and set of variable values from the user, and passes these to
// lpo function. It then displays the results obtained and any errors that may have
// occurred. Function accepts no input. In case of failure, it returns an error.
func wpCalcLhs() error {
	var index      int      // index of row we want to check
	var status     int      // status of operation returned back to us
	var lhs        float64  // lhs value of the constraint
	var conPoint []float64  // set of variable values for the constraint
	var err        error    // error returned by functions called

	// Get the row index and point from the user, and pass these to the lpo
	// function provided there were no errors.
	if err = wpGetPoint(&index, &conPoint); err != nil {
		return errors.Wrap(err, "wpCalcLhs failed to get point")
	} 

	if err = lpo.CalcLhs(index, conPoint, &lhs, &status); err != nil {
		return errors.Wrap(err, "wpCalcLhs failed to get lhs")
	}

	fmt.Printf("\nLHS = %f, Status = %d.\n", lhs, status)
	
	return nil	
}

//==============================================================================

// wpCalcConViol is a wrapper for lpo.CalcConViolation. It uses a separate function 
// to get the constraint index and set of variable values from the user, and passes 
// these to the lpo function. It then displays the results obtained and any errors 
// that may have occurred. Function accepts no input. 
// In case of failure, it returns an error.
func wpCalcConViol() error {
	var index      int      // index of row we want to check
	var status     int      // status of operation that was performed
	var conViol    float64  // constraint violation magnitude and sign
	var conPoint []float64  // set of variable values for the constraint
	var err        error    // error returned by functions called

	// Get the row index and point from the user, and pass these to the lpo
	// function provided there were no errors.

	if err = wpGetPoint(&index, &conPoint); err != nil {
		return errors.Wrap(err, "wpCalcConViol failed to get point")
	} 
	
	err = lpo.CalcConViolation(index, conPoint, &conViol, &status)
	fmt.Printf("Violation = %f, Status = %d\n", conViol, status)
	
	if err != nil {
		errors.Wrap(err, "wpCalcConViol failed to get violation")
	}
	
	return nil		
}

//==============================================================================

// runMainWrapper executes the "main" wrapper that is available for external users
// rather than the development wrapper.
// The function is called from the highest-level wrapper and accepts the cmdOption as an 
// argument. If the command cannot be executed because it does not match any of 
// the cases covered by this wrapper, it returns an error.
func runMainWrapper(cmdOption string) error {
	var fileName      string        // name of file being processed
//	var userString    string        // holder for string input by user
//	var userInt       int           // holder for int input by user
//	var tmpString     string        // temp holder for string variables
//	var tmpInt        int           // temp holder for int variables
	var err           error         // error returned by functions called

	switch cmdOption {

		case "1":
			// Read MPS file
			fmt.Printf("Enter name of MPS file to be read: ")
			fmt.Scanln(&fileName)
			if custEnvOn {
				fileName = dSrcDev + fileName + fExtension
			}
			fmt.Println("Reading file", fileName)
			if err = lpo.ReadMpsFile(fileName); err != nil {
				fmt.Println(err)
			}

		case "2":
			// Write MPS file
			fmt.Printf("Enter MPS output file name: ")
			fmt.Scanln(&fileName)
			if custEnvOn {
				fileName = dSrcDev + fileName + fExtension
			}
			err = lpo.WriteMpsFile(fileName)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Model successfully written to file '%s'.\n", fileName)
			}

		case "3":
			// Solve problem
			err = wpSolveProb(0)
			if err != nil {
				fmt.Println(err)
			}

		case "4":
			// ReduceMatrix
			err = wpReduceMtrx()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("\nExample showing ReduceMatrix completed successfully.\n")
			}

		case "5":
			wpInitLpo()	
			
		case "6":
			// Print LPO input data structures
			wpPrintLpoIn()
			
		case "7":
			// Print LPO solution data structures
			wpPrintLpoSoln()	

		case "8":
			// Print Cplex solution
			wpPrintCplexSoln()

		//----------------------------------------------------------------------
		default:
			return errors.Errorf("Command %s not in user menu", cmdOption)
		
	} // end switch on command option

	return nil
			
}

//==============================================================================

// runLpoWrapper executes functions from the LPO package. 
// The display of menu items may be hidden to avoid clutter, but the command
// options remain available even if the menu item is hidden. 
// The function is called from the main wrapper and accepts the cmdOption as an 
// argument. If the command cannot be executed because it does not match any of 
// the cases covered by this wrapper, it returns an error.
func runLpoWrapper(cmdOption string) error {	
	var fileName      string        // name of file being processed
	var userString    string        // holder for string input by user
	var userInt       int           // holder for int input by user
	var tmpString     string        // temp holder for string variables
	var tmpInt        int           // temp holder for int variables
	var err           error         // error returned by functions called

	// The gpx variables used in this function are package global variables so
	// we don't have to pass them to the higher-level wrapper and back again as
	// individual commands that use them are executed.
	
	switch cmdOption {

	//--------------------------------------------------------------------------
	case "21":
		fmt.Printf("Adjusting model.\n")
		if err = lpo.AdjustModel(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Post-processing on model completed successfully.\n")				
		}			

	//--------------------------------------------------------------------------
	case "22":
		// CalcConViolation
		if err = wpCalcConViol(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("CalcConViolation completed successfully.\n")
		}

	//--------------------------------------------------------------------------
	case "23":
		// CalcLhs
		if err = wpCalcLhs(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("CalcLhs completed successfully.\n")
		}		

	//--------------------------------------------------------------------------
	case "24":
		// CoinParseSoln
		fmt.Printf("Command not yet implemented.\n")

	//--------------------------------------------------------------------------
	case "25":
		// CoinSolveMps
		fmt.Printf("Command not yet implemented.\n")

	//--------------------------------------------------------------------------
	case "26":
		// CoinSolveProb
		err = wpSolveProb(2)
		if err != nil {
			fmt.Println(err)
		}

	// 27: CplexCreateProb in utilsgpx
	//--------------------------------------------------------------------------
	case "28":
		fmt.Printf("\nEnter file name containing Cplex output: ")
		fmt.Scanln(&userString)
		if custEnvOn {
			userString = dSrcDev + userString + fExtension
		}
		if err = lpo.CplexParseSoln(userString, &lpCpSoln); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("CplexParseSoln completed successfully.\n")
		}

	//--------------------------------------------------------------------------
	case "29":
		// Read and solve MPS file directly by Cplex
		err = wpSolveCplex()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("\nExample using Cplex directly completed successfully.\n")
		}

	//--------------------------------------------------------------------------
	case "30":
		// CplexSolveProb
		err = wpSolveProb(1)
		if err != nil {
			fmt.Println(err)
		}

	//--------------------------------------------------------------------------
	case "31":
		fmt.Printf("Enter index of column to delete: ")
		fmt.Scanln(&userInt)
		if err = lpo.DelCol(userInt); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Column %d successfully deleted.\n", userInt)
		}

	//--------------------------------------------------------------------------
	case "32":
		fmt.Printf("Enter index of row to delete: ")
		fmt.Scanln(&userInt)
		if err = lpo.DelRow(userInt); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Row %d successfully deleted.\n", userInt)
		}

	//--------------------------------------------------------------------------
	case "33":
		if err = lpo.GetLogLevel(&tmpInt); err != nil {
			fmt.Println(err)				
		} else {
			fmt.Printf("Log level is set to %d.\n", tmpInt)
		}

	//--------------------------------------------------------------------------
	case "34":
		if err = lpo.GetStatistics(&lpStats); err != nil {
			fmt.Println(err)				
		} else {
			fmt.Printf("Statistics successfully obtained.\n")
		}

	//--------------------------------------------------------------------------
	case "35":
		if err = lpo.GetTempDirPath(&tmpString); err != nil {
			fmt.Println(err)				
		} else {
			fmt.Printf("Temp dir set to %s.\n", tmpString)
		}
		
	//--------------------------------------------------------------------------
	case "36":
		if err = lpo.InitModel(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Model successfully initialized.\n")
		}

	//--------------------------------------------------------------------------
	case "37":
		fmt.Printf("Enter index of column to print: ")
		fmt.Scanln(&userInt)
		if err = lpo.PrintCol(userInt); err != nil {
			fmt.Println(err)				
		}

	//--------------------------------------------------------------------------
	case "38":
		if err = lpo.PrintModel(); err != nil {
			fmt.Println(err)
		}
		
	//--------------------------------------------------------------------------
	case "39":
		if err = lpo.PrintRhs(); err != nil {
			fmt.Println(err)
		}

	//--------------------------------------------------------------------------
	case "40":
		fmt.Printf("Enter index of row to print: ")
		fmt.Scanln(&userInt)
		if err = lpo.PrintRow(userInt); err != nil {
			fmt.Println(err)				
		}

	//--------------------------------------------------------------------------
	case "41":
		// PrintStatistics
		if lpStats.NumElements == 0 {
			fmt.Printf("There are no elements in data structure.\n")
		} else {
			if err = lpo.PrintStatistics(lpStats); err != nil {
				fmt.Println(err)
			}			
		}

	//--------------------------------------------------------------------------
	case "42":
		// Read MPS file
		fmt.Printf("Enter name of MPS file to be read: ")
		fmt.Scanln(&fileName)
		if custEnvOn {
			fileName = dSrcDev + fileName + fExtension
		}
		fmt.Println("Reading file", fileName)
		if err = lpo.ReadMpsFile(fileName); err != nil {
			fmt.Println(err)
		}

	//--------------------------------------------------------------------------
	case "43":
		// ReduceMatrix
		err = wpReduceMtrx()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("\nExample showing ReduceMatrix completed successfully.\n")
		}

	//--------------------------------------------------------------------------
	case "44":
		if err = lpo.ScaleRows(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Rows scaled successfully.\n")
		}

	//--------------------------------------------------------------------------
	case "45":
		fmt.Printf("Enter new log level: ")
		fmt.Scanln(&userInt)
		if err = lpo.SetLogLevel(userInt); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Log level changed to %d.\n", userInt)
		}

	//--------------------------------------------------------------------------
	case "46":
		userString = ""
		fmt.Printf("Enter new path for temp directory: ")
		fmt.Scanln(&userString)
		if err = lpo.SetTempDirPath(userString); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Temp dir changed to %s.\n", userString)
		}

	//--------------------------------------------------------------------------
	case "47":
		fmt.Printf("Enter number of TightenBounds iterations: ")
		fmt.Scanln(&userInt)
		if err = lpo.TightenBounds(userInt, &tmpInt); err != nil {
			fmt.Println(err)								
		}
		fmt.Printf("TightenBounds completed %d of %d iterations\n", userInt, tmpInt)

	// 48: TransFromGpx in utilsgpx
	// 49: TransToGpx   in utilsgpx
	//--------------------------------------------------------------------------
	case "50":
		// Write MPS file
		fmt.Printf("Enter MPS output file name: ")
		fmt.Scanln(&fileName)
		if custEnvOn {
			fileName = dSrcDev + fileName + fExtension
		}
		err = lpo.WriteMpsFile(fileName)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Model successfully written to file '%s'.\n", fileName)
		}

	//--------------------------------------------------------------------------
	case "51":
		userString = ""
		fmt.Printf("Enter name of PSOP file: ")
		fmt.Scanln(&userString)
		fmt.Printf("Enter number of coef per line, <0 for all, 0 for none: ")
		fmt.Scanln(&userInt)				
							
		if custEnvOn {
			tmpString  = dSrcDev + fPrefPsopOut + userString + fExtension			
		} 
						
		if err = lpo.WritePsopFile(tmpString, userInt); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("PSOP written to file '%s'\n.", tmpString)
		}

	//--------------------------------------------------------------------------
	default:
		return errors.Errorf("Command %s not in functions menu", cmdOption)
		
	} // End switch on command option
	
	return nil	
}