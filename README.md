# Command runopt

The runopt executable is a wrapper for exercising functions exported by the lpo and gpx packages which provide a suite
of Go language tools for Linear Programming (LP) and Mixed-Integer Linear Programming (MILP). It was developed as a test
tool and assumes that the user is aware of the proper sequence in which individual functions need to be called in order
to achieve the desired result.

# Dependencies

The runopt executable is dependent on the following:

*	github.com/pkg/errors
* github.com/go-opt/lpo
*	github.com/go-opt/gpx (if using the callable C functions provided by Cplex)

The lpo and gpx packages are themselves dependent on the installation and configuration of the Cplex solver and a
C compiler or the Coin-OR solver. Please refer to those packages for details.


# Installation and Configuration

To install the executable on a Windows platform, go to the cmd.exe window and enter the command:
```
  go get -u github.com/go-opt/runopt
```

The default configuration of runopt assumes that both Cplex and Coin-OR are installed. If Coin-OR is not installed,
no modifications to the default configuration are needed. The only impact in such a case is that functions testing
Coin-OR functionality will return an error. However, if Cplex is not installed, the default configuration must be
changed to avoid compilation failures.

## Configuring runopt without gpx

If gpx is not installed, you need to modify the following files so that the other functions not using gpx may be 
compiled and executed.

File utilsgpx.go must be excluded from being built by uncommenting the first line of that file so that it reads:
```
  // +build exclude
```
In the runopt.go file, search for the GPX_EXCLUDED string and comment out the command which immediately follows
so that the new code reads as follows:
```
  ...

  err = errors.New("Requested solver not present")
  if useCoinSolver {
    err = lpo.CoinSolveProb(psCtrl, &psResult)						
  } else {
    // GPX_EXCLUDED: Comment out the following line if gpx is not installed.
    // err = lpo.CplexSolveProb(psCtrl, &psResult)			
  }
  
  ...
  
  if gpxMenuOn {
    err = errors.New("Command not available")
    // GPX_EXCLUDED: Comment out the following line if gpx is not installed.
    // err = runGpxWrapper(cmdOption)
    if err == nil {
      // Found the command in gpx menu
      continue
    }
  }  
  ...
```
Once you have made these changes, you can compile and use runopt without gpx.
