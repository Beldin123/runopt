// 01 - Jul. 12, 2018   First version, uploaded to github
// 02 - Sept. 6, 2018   Revised as first "delux" lporun now replacing old runopt


/* 

Test tool for exercising lpo and gpx functions.

SUMMARY


This executable serves as a wrapper to allow users to populate key input data,
inspect the contents of data structures, and call functions exported by the lpo
and gpx packages. 

Commands in the main menu exercise the most important functionality, or contain 
commands which don't correspond to any exported function. Exported functions 
are listed in alphabetical order in menus that must be explicitly toggled by 
pressing "s" for lpo and "g" for gpx.

The options available from the main menu are:

    0 - exit program
    1 - read MPS file (needed by other functions called at a later time)
    2 - write MPS file (useful if user-written functions populated lpo structures)
    3 - solve problem (read MPS if model not already loaded, reduce, and solve)
    4 - reduce matrix to reduce but not solve a problem
    5 - initialize lpo structures (needed in conjunction with lpo function exerciser)
    6 - show lpo input data structures
    7 - show lpo solution provided by the solver
    8 - show Cplex solution (contents of xml file loaded into data structures)
    9 - initialize gpx structures (needed in conjunction with gpx function exerciser)
   10 - show gpx input data structures
   11 - show gpx solution data structures


Toggles control the following functionality:	

    c - toggle custom environment (to reduce typing when entering file names)
    s - toggle lpo function exerciser (to enable access to exported lpo functions)
    g - toggle gpx function exerciser (to enable access to exported gpx functions)

To select an option, enter the corresponding letter or number when prompted.


MAIN COMMANDS

The main command options are always enabled and no option has been provided to 
disable them.


Exit

This option is used to terminate execution of the program. This option is displayed
as part of the command prompt and is not included in the lists showing other options.


Read MPS file

This option uses the ReadMpsFile function to populate the internal lpo data structures
from an MPS file. Although this single function is included in the lpo function
exerciser, it is important enough to be included in the main menu.


Write MPS file

Similarly, this option consists of the WriteMpsFile function which is also considered
important enough to be included in the main menu.


Solve problem

This option is used to load a model into lpo (or use the model loaded by a
previous command), reduce the problem size, and solve it via Coin-OR or Cplex.
All of the work is done by the CoinSolveProb or CplexSolveProb function, based on
responses to prompts provided by the user in order.

The first prompt the user must answer is the source from which the model is to be
read. This can be the name of an MPS file, or can be left empty (carriage return at
the prompt) if the model is already loaded into the internal data structures as
a result of an earlier operation, and should be taken from there.

The next set of prompts allows the user to specify the file names for storing
the Cplex solution (xml), reduced matrix (MPS), and pre-solve operations (text).
These files are optional, and if the name is not entered, the file will not be written.
If the custom environment is enabled, the prompt for these file names will not
appear, it will be assumed that all of them are needed, and the name will be based
on the file name provided with the appropriate prefix added (see custom environment
section for details).

The next prompt allows the user to set the solver to be used, either Coin-OR or Cplex.

The next prompt allows the user to specify which matrix-reduction operations to
apply, and whether to solve the problem. The high-level options are "all" (apply all
reductions and solve the problem), "none" (don't reduce anything but solve problem),
or blank (carriage-return) to set each flag independently. The user must explicitly
enter "Y" at each of the prompts to set the corresponding to "true", since any other
response to the prompt will leave the flag in its default "false" state.

After all prompts have been answered, the populated control structure is passed
to CoinSolveProb or CplexSolveProb which returns a solution, or an error. If no
error occurred, the user has the option to display the results. The results may
also be displayed at a later time using the "Show lpo solution" option.


Reduce matrix

This example is a subset of the "Using SolveProb" example. The model must be loaded
into the internal data structures, most likely by an earlier call to ReadMpsFile,
TransFromGpx if converting from gpx data structures, or some other similar
mechanism. The user is prompted to specify which matrix-reduction operations are
to be performed using the same set of questions requiring "Y" or "N" responses,
the problem is reduced, and the system is left in this state. The user may then
perform additional operations by independently calling other lpo or gpx functions
as needed.

Initialize lpo structures

This option initializes all data structures used in this program. It is more
thorough than InitModel, which only initializes the input data structures. This
option is intended to be used in conjunction with the lpo and/or gpx function
exercisers.

Show lpo input

This option shows the lpo input data structures in their raw form. It is not "pretty",
but displays all fields of the various lists, and is useful when exercising other
functions (e.g. DelRow or DelCol). To display a prettier version of the model,
please use one of the other "Print" functions provided for this purpose.

Show lpo solution

This option shows the lpo solution data structures. It is intended to be used
in conjunction with the function exerciser.

Show Cplex solution

This option displays the data structure containing the solution obtained by parsing
the Cplex solution xml file. It is useful when wishing to look at the raw Cplex
solution without having to open the file.

Initialize gpx structures

The gpx package is a key component of lpo, and a function exerciser has been provided
for this package as well. This option is intended to initialize the gpx data structures
so that the functions in both packages can be used.

Show gpx input

This option shows the gpx input data structures. It is useful when exercising
the gpx component, which acts as the interface between lpo and Cplex.

Show gpx solution

This option is used to display the solution provided by Cplex. It is useful when
running individual gpx functions which do not automatically show the solution when
it is obtained.


TOGGLES

This section describes the toggles which control program behaviour. The variables
which control the toggles and their default state are:

   var lpoMenuOn  bool = false   // Flag for enabling lpo functions   
   var gpxMenuOn  bool = false   // Flag for enabling gpx functions   
   var custEnvOn  bool = false   // Flag for enabling custom paths and names


Toggle lpo function exerciser

This toggle is used to enable or disable the options which are available to exercise
individual lpo functions. By default, the lpo function exerciser is disabled.

Toggle gpx function exerciser

This toggle is used to enable or disable the options which are available to exercise
individual gpx functions. By default, the gpx function exerciser is disabled.

Toggle for custom environment

This toggle controls how file names are handled by this program. If all files are
located in the same directory and if all files have the same extension, this
option reduces the amount of typing needed to answer various prompts. By default,
the custom environment is disabled and the full file name (including path and
extension) must be specified.

If custom environment is enabled (variable set to "true"), the directory name is
added as a prefix to the base file name, the extension is added as a suffix
to that name, and any "family" of files (e.g. Cplex output, PSOP file, etc.) is based
on the core name input by the user but with a prefix added to it. The default
settings are:

  var dSrcDev       string = "D:/Docs/LP/Data/"           // Development source data dir
  var fPrefSolnOut  string = "sol_"   // Prefix for solution xml files  
  var fPrefRdcMps   string = "rmx_"   // Prefix for MPS file storing reduced matrix
  var fPrefPsopOut  string = "psop_"  // Prefix for file storing data removed during PSOP
  var fExtension    string = ".txt"   // Extension of source data files in development dir.  

Caution is advised if using a custom environment.

LPO FUNCTION EXERCISER

This section lists the options used to exercise individual gpx functions. Please
refer to the main documentation for details on function input, output, and behaviour.

Care must be taken that the required data structures have been correctly initialized 
and populated, and that the functions are not called out of sequence. The list of 
available functions, listed in alphabetical order, and some things to watch out for, 
are listed below.

 21 - AdjustModel      - Do post-processing after data structures are populated.
 22 - CalcConViolation - Calculate the constraint violation for a given point.
 23 - CalcLhs          - Calculate the LHS for a given point.
 24 - CoinParseSoln    - Parse the xml solution file generated by Coin-OR.
 25 - CoinSolveMps     - Have Coin-OR solve the problem defined in the MPS file.
 26 - CoinSolveProb    - Reduces and solves the model via the Coin-OR solver.
 27 - CplexCreateProb  - Initialize Cplex environment and convert to gpx structures.
 28 - CplexParseSoln   - Parse Cplex xml solution file into internal structures.
 29 - CplexSolveMps    - Have Cplex solve the problem defined in the MPS file.
 30 - CplexSolveProb   - Reduces and solves the model via Cplex callable libraries.
 31 - DelCol           - Delete a specific column from the lpo columns list.
 32 - DelRow           - Delete a specific row from the lpo rows list.
 33 - GetLogLevel      - Get the current log level.
 34 - GetStatistics    - Get the model statistics.
 35 - GetTempDirPath   - Get the current path of the temp directory.
 36 - InitModel        - Initialize the lpo input data structures.
 37 - PrintCol         - Prints the rows in which the column, specified by its index, occurs.
 38 - PrintModel       - Prints the model in equation format.
 39 - PrintRhs         - Prints the RHS of all constraints.
 40 - PrintRow         - Prints the row, specified by its index, in equation format.
 41 - PrintStatistics  - Prints the model statistics.
 42 - ReadMpsFile      - Reads MPS file and populates internal data structures.
 43 - ReduceMatrix     - Performs the matrix-reduction operations specified.
 44 - ScaleRows        - Performs row scaling on the entire model.
 45 - SetLogLevel      - Sets the log level to the value specified.
 46 - SetTempDirPath   - Sets the temp dir location to the path specified.
 47 - TightenBounds    - Tightens the bounds on the constraints of the model.
 48 - TransFromGpx     - Populates lpo data structures from the gpx data structures.
 49 - TransToGpx       - Populates gpx data structures from the lpo data structures.
 50 - WriteMpsFile     - Writes the model to an MPS file.
 51 - WritePsopFile    - Writes the pre-solve operations (PSOP) to a text file.

GPX FUNCTION EXERCISER

This section lists the options used to exercise individual gpx functions. The same
functionality is available in an executable included with the gpx package, but has
been included here for convenience. Please refer to the main documentation for 
details on function input, output, and behaviour.

Care must be taken that the required data structures have been correctly initialized 
and populated, and that the functions are not called out of sequence. The list of 
available functions, listed in alphabetical order, and some things to watch out for, 
are listed below.

 61 - ChgCoefList     - Sets non-zero coefficients, must be used after NewCols and NewRows.
 62 - ChgObjSen       - Sets problem to be treated as "maximize" or "minimize".
 63 - ChgProbName     - Sets the problem name.
 64 - CloseCplex      - Cleans up and closed the Cplex environment, must be called last.
 65 - CreateProb      - Initializes the Cplex environment, must be called first.
 66 - GetColName      - Creates the column solution list of the correct size and 
                        populates it with the column names.
 67 - GetMipSolution  - Creates and populates solution structures for MIP problem.
 68 - GetNumCols      - Gets the number of columns in the problem.
 69 - GetNumRows      - Gets the number of rows in the problem.
 70 - GetObjVal       - Gets the obj. func. value, assumes problem has been solved.
 71 - GetRowName      - Creates the row solution list of the correct size and populates 
                        it with the row names.
 72 - GetSlack        - Adds slack values to row solution list, which must exist.
 73 - GetSolution     - Creates and populates solution structures for LP problem.
 74 - GetX            - Adds values to column solution list, which must exist.
 75 - LpOpt           - Optimizes an LP loaded into Cplex.
 76 - MipOpt          - Optimizes a MIP loaded into Cplex.
 77 - NewCols         - Creates new columns in Cplex from the internal data structures.
 78 - NewRows         - Creates new rows in Cplex from the internal data structures.
 79 - OutputToScreen  - Specifies whether Cplex should display output to screen or not.
 80 - ReadCopyProb    - Populates the problem in Cplex directly from the file specified.
 81 - SolWrite        - Writes the Cplex solution to a file.
 82 - WriteProb       - Writes the problem loaded into Cplex to a file using the
                        format specified.



*/
package main