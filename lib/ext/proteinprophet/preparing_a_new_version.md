1. Add Boost libraries to the header

```cpp
#include "pwiz/utility/misc/random_access_compressed_ifstream.hpp"

// line 76 / 77
#include <boost/filesystem/operations.hpp>
#include <boost/filesystem/path.hpp>

typedef std::map<ppProtein *, ppProtein *> ProteinProteinMap;
```

1.  Add the alternative directory to the forward declarations

```cpp
std::string BINARY_DIRECTORY;
std::string SERVER_ROOT;
// line 292
std::string ALT_BINARY_DIRECTORY;
std::string WEBSERVER_URL;
std::string TPPVersionInfo;
```

1. Add selfpath to ALT bin directory

```cpp
// abort if fewer than two command line arguments
if ( argc < 3 ) {
  usage();
}

namespace fs = boost::filesystem;
fs::path selfpath = argv[0];
// line 773
ALT_BINARY_DIRECTORY = selfpath.remove_filename().string() + "/";
```

1. Replace the bin path for binaries with ALT path

```cpp
//std::string quanticparser = BINARY_DIRECTORY + "QuanticProteinParser";
//std::string xpressparser = BINARY_DIRECTORY + "XPressProteinRatioParser";
//std::string asapparser = BINARY_DIRECTORY + "ASAPRatioProteinRatioParser";
//std::string asappvalueparser = BINARY_DIRECTORY + "ASAPRatioPvalueParser";
//std::string libraparser = BINARY_DIRECTORY + "LibraProteinRatioParser";

// line 1058
std::string quanticparser = ALT_BINARY_DIRECTORY + "QuanticProteinParser";
std::string xpressparser = ALT_BINARY_DIRECTORY + "XPressProteinRatioParser";
std::string asapparser = ALT_BINARY_DIRECTORY + "ASAPRatioProteinRatioParser";
std::string asappvalueparser = ALT_BINARY_DIRECTORY + "ASAPRatioPvalueParser";
std::string libraparser = ALT_BINARY_DIRECTORY + "LibraProteinRatioParser";

// line 1285
// first get database
//std::string database_parser = BINARY_DIRECTORY + "DatabaseParser";
//std::string refresh_parser = BINARY_DIRECTORY + "RefreshParser";

std::string database_parser = ALT_BINARY_DIRECTORY + "DatabaseParser";
std::string refresh_parser = ALT_BINARY_DIRECTORY + "RefreshParser";

group_weight
```

1. Update the computeCoverage function

```cpp
//##############################################################################
// function computeCoverage
//    input:   double - minimum protein probability
//    input:   int - maximum number of proteins
//    output:  none
//##############################################################################
void computeCoverage( double min_prot_prob, int max_num_prots )
{
std::string covinfofile = OUTFILE + ".covinfo";
std::string covresultsfile = OUTFILE + ".cov";
//std::string coverage_exec = BINARY_DIRECTORY + "batchcoverage";
std::string coverage_exec = ALT_BINARY_DIRECTORY + "batchcoverage";
if (!database.empty() && !fileExists( database ) ) {  // look in the wwwroot
  std::string strtry = resolve_root(database.c_str());
  if ( fileExists(strtry)) {
     database = strtry;
  }
}
```

1. Add group weight back to output

```cpp
// line 3319
outfile << "group_weight=\"" << std::setprecision( 2 ) << parentProt.getGroupPepWt() << "\" ";
```