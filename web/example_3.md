# Processing, Filtering and Analyzing a CPTAC 3 cohort

For this example we will see how to process and analyzed the Clear Cell Renal Carcinoma cohort data from CPTAC 3 using MSFragger and Philosopher. You will learn how to process a large cohort composed by multiple fractionated TMT-labeled data sets. This tutorial will contain as much details as possible for you to reproduce our results on your side. Keep in mind that it is expected to see small differences between some results once that these tools are under constant improvement.

We will need:

* [Philosopher](https://prvst.github.io/philosopher/) (version 20181218 or higher).
* [MSFragger](https://www.nature.com/articles/nmeth.4256) (version 20180316 or higher).
* [Java](http://www.oracle.com/technetwork/java/javase/downloads/jre9-downloads-3848532.html) version 9 (MSFragger requirement).
* The Clear Cell Renal Carcinoma data set from [CPTAC 3](https://cptacdcc.georgetown.edu/cptac/).
* A computer server running GNU/Linux with at least 64GB of RAM.

We ran this example on a Linux Red Hat 7, meaning that the commands you see below will be "Linux compatible", if you are trying to reproduce this on a Windows machine, you will need to adjust the folder separators ('\\' for windows and '/' for Linux).


## Downloading the data set
The CPTAC 3 data is currently available at the NIH / CPTAC Private Data portal, if you are not part of the consortium you will need to sign an agreement that prohibits any publication using this data until the embargo has concluded. We will not need to convert the RAW files because we are using the mzML provided by the consortium.


## Downloading our tools
You can download all the necessary software using the links provided above, make sure to always get the latest version available.


## Organizing our directory
Having all files in an organized way is important for later when running our pipeline, we start by creating a folder for the entire cohort data analysis that will be called _6_CPTAC3_Clear_Cell_Renal_Cell_Carcinoma_, inside we will create one folder for the whole proteome analysis and another for the phosphoproteome, then inside each, individual folders for the TMT fractionated experiment (here called data set), this is what the directory looks like:

```
6_CPTAC3_Clear_Cell_Renal_Cell_Carcinoma
├── phospho
│   ├── 01CPTAC_CCRCC_P_JHU_20171106
│   ├── 02CPTAC_CCRCC_P_JHU_20171108
│   ├── 03CPTAC_CCRCC_P_JHU_20171110
│   ├── 04CPTAC_CCRCC_P_JHU_20171113
│   ├── 05CPTAC_CCRCC_P_JHU_20171117
│   ├── 06CPTAC_CCRCC_P_JHU_20171124
│   ├── 07CPTAC_CCRCC_P_JHU_20171212
│   ├── 08CPTAC_CCRCC_P_JHU_20171214
│   ├── 09CPTAC_CCRCC_P_JHU_20180130
│   ├── 10CPTAC_CCRCC_P_JHU_20180123
│   ├── 11CPTAC_CCRCC_P_JHU_20180201
│   ├── 12CPTAC_CCRCC_P_JHU_20180206
│   ├── 13CPTAC_CCRCC_P_JHU_20180221
│   ├── 14CPTAC_CCRCC_P_JHU_20180226
│   ├── 15CPTAC_CCRCC_P_JHU_20180321
│   ├── 16CPTAC_CCRCC_P_JHU_20180326
│   ├── 17CPTAC_CCRCC_P_JHU_20180411
│   ├── 18CPTAC_CCRCC_P_JHU_20180524
│   ├── 19CPTAC_CCRCC_P_JHU_20180529
│   ├── 20CPTAC_CCRCC_P_JHU_20180531
│   ├── 20CPTAC_CCRCC_P_JHU_20180619
│   ├── 21CPTAC_CCRCC_P_JHU_20180613
│   ├── 22CPTAC_CCRCC_P_JHU_20180615
│   ├── 23CPTAC_CCRCC_P_JHU_20180617
|── whole
|   ├── 01CPTAC_CCRCC_W_JHU_20171007
|   ├── 02CPTAC_CCRCC_W_JHU_20171003
|   ├── 03CPTAC_CCRCC_W_JHU_20171022
|   ├── 04CPTAC_CCRCC_W_JHU_20171026
|   ├── 05CPTAC_CCRCC_W_JHU_20171030
|   ├── 06CPTAC_CCRCC_W_JHU_20171120
|   ├── 07CPTAC_CCRCC_W_JHU_20171127
|   ├── 08CPTAC_CCRCC_W_JHU_20171205
|   ├── 09CPTAC_CCRCC_W_JHU_20171215
|   ├── 10CPTAC_CCRCC_W_JHU_20180119
|   ├── 11CPTAC_CCRCC_W_JHU_20180126
|   ├── 12CPTAC_CCRCC_W_JHU_20180202
|   ├── 13CPTAC_CCRCC_W_JHU_20180215
|   ├── 14CPTAC_CCRCC_W_JHU_20180223
|   ├── 15CPTAC_CCRCC_W_JHU_20180315
|   ├── 16CPTAC_CCRCC_W_JHU_20180322
|   ├── 17CPTAC_CCRCC_W_JHU_20180517
|   ├── 18CPTAC_CCRCC_W_JHU_20180521
|   ├── 19CPTAC_CCRCC_W_JHU_20180526
|   ├── 20CPTAC_CCRCC_W_JHU_20180602
|   ├── 21CPTAC_CCRCC_W_JHU_20180621
|   ├── 22CPTAC_CCRCC_W_JHU_20180625
|   ├── 23CPTAC_CCRCC_W_JHU_20180629
├── bin
│   ├── MSFragger-20180316.jar
│   ├── philosopher
|── params
|   ├── fragger.params
|   ├── philosopher.yaml
|── database
|   └── 2018-07-16-td-RefSeq.20180629_Human_ucsc_hg38_cpdbnr_mito_264contams.fasta
```

Note that I also created a folder called `bin` for the software we are going to use, a folder called `params` for the parameters files and a folder called `database` with our protein FASTA database.

For the sake of simplicity, I will only show how to process the _whole proteome_.

Inside each one of these folders we will place the mzML files corresponding to the given fractions and one annotation file for the TMT channels:

```
.
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f01.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f02.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f03.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f04.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f05.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f06.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f07.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f08.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f09.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f10.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f11.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f12.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f13.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f14.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f15.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f16.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f17.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f18.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f19.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f20.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f21.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f22.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f23.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_f24.mzML
├── 01CPTAC_CCRCC_W_JHU_20171007_LUMOS_fA.mzML
├── annotation.txt
```

The annotation file is a simple text file with a map between the TMT channels and the sample labels, this will be useful at the end when we have the final report. Each data set folder should contain a text file called _annotation.txt_ with the mapping. Below is an example of the annotation file for the data set #01:

```
126 CPT0079430001
127N CPT0023360001
127C CPT0023350003
128N CPT0079410003
128C CPT0087040003
129N CPT0077310003
129C CPT0077320001
130N CPT0087050003
130C CPT0002270011
131N pool01
```

The given labels for each cohort and data set can also be found on the NIH CPTAC data portal.


## Setting the MSFragger parameter file

We are going to use the parameter file displayed below for our analysis. you can find more details about each parameter on the MSFragger manual.

```
database_name = /workspace/6_CPTAC3_Clear_Cell_Renal_Cell_Carcinoma/database/RefSeq_20180629/2018-07-16-td-RefSeq.20180629_Human_ucsc_hg38_cpdbnr_mito_264contams.fasta
num_threads = 28                        # 0=poll CPU to set num threads; else specify num threads directly (max 64)

precursor_mass_tolerance = 20.0
precursor_mass_units = 1               # 0=Daltons, 1=ppm

precursor_true_tolerance = 20.00
precursor_true_units = 1
fragment_mass_tolerance = 20.00
fragment_mass_units = 1		             # 0=Daltons, 1=ppm


isotope_error = -1/0/1/2/3             # 0=off, -1/0/1/2/3 (standard C13 error)

search_enzyme_name = Trypsin
search_enzyme_cutafter = KR
search_enzyme_butnotafter = P

num_enzyme_termini = 2                 # 2 for enzymatic, 1 for semi-enzymatic, 0 for nonspecific digestion
allowed_missed_cleavage = 2            # maximum value is 5

clip_nTerm_M = 1

#maximum of 7 mods - amino acid codes, * for any amino acid, [ and ] specifies protein termini, n and c specifies peptide termini
variable_mod_01 = 15.9949 M
variable_mod_02 = 42.0106 [^
variable_mod_03 = 229.162932 n^
variable_mod_04 = 229.162932 S

allow_multiple_variable_mods_on_residue = 1  	# static mods are not considered
max_variable_mods_per_mod = 3 			          # maximum of 5
max_variable_mods_combinations = 50000  	    # maximum of 65534, limits number of modified peptides generated from sequence

output_format = pepXML
output_file_extension = pepXML   #pep.xml
output_report_topN = 3
output_max_expect = 50

precursor_charge = 0 0                 # precursor charge range to analyze; does not override any existing charge; 0 as 1st entry ignores parameter
override_charge = 0                    # 0=no, 1=yes to override existing precursor charge states with precursor_charge parameter
ms_level = 2                           # MS level to analyze, valid are levels 2 (default) or 3

digest_min_length = 7
digest_max_length = 50
digest_mass_range = 500.0 5000.0       # MH+ peptide mass range to analyze
max_fragment_charge = 2                # set maximum fragment charge state to analyze (allowed max 5)

#open search parameters
track_zero_topN = 0		             # in addition to topN results, keep track of top results in zero bin
zero_bin_accept_expect = 0	       # boost top zero bin entry to top if it has expect under 0.01 - set to 0 to disable
zero_bin_mult_expect = 1	         # disabled if above passes - multiply expect of zero bin for ordering purposes (does not affect reported expect)
add_topN_complementary = 0

# spectral processing

minimum_peaks = 15                     # required minimum number of peaks in spectrum to search (default 10)
use_topN_peaks = 100
min_fragments_modelling = 3
min_matched_fragments = 4
minimum_ratio = 0.01		                 # filter peaks below this fraction of strongest peak
clear_mz_range = 125.5 131.5             # for iTRAQ/TMT type data; will clear out all peaks in the specified m/z range

# additional modifications

add_Cterm_peptide = 0.0
add_Nterm_peptide = 0.0
add_Cterm_protein = 0.0
add_Nterm_protein = 0.0

add_G_glycine = 0.0000                 # added to G - avg.  57.0513, mono.  57.02146
add_A_alanine = 0.0000                 # added to A - avg.  71.0779, mono.  71.03711
add_S_serine = 0.0000                  # added to S - avg.  87.0773, mono.  87.03203
add_P_proline = 0.0000                 # added to P - avg.  97.1152, mono.  97.05276
add_V_valine = 0.0000                  # added to V - avg.  99.1311, mono.  99.06841
add_T_threonine = 0.0000               # added to T - avg. 101.1038, mono. 101.04768
add_C_cysteine = 57.021464             # added to C - avg. 103.1429, mono. 103.00918
add_L_leucine = 0.0000                 # added to L - avg. 113.1576, mono. 113.08406
add_I_isoleucine = 0.0000              # added to I - avg. 113.1576, mono. 113.08406
add_N_asparagine = 0.0000              # added to N - avg. 114.1026, mono. 114.04293
add_D_aspartic_acid = 0.0000           # added to D - avg. 115.0874, mono. 115.02694
add_Q_glutamine = 0.0000               # added to Q - avg. 128.1292, mono. 128.05858
add_K_lysine = 229.162932              # added to K - avg. 128.1723, mono. 128.09496
add_E_glutamic_acid = 0.0000           # added to E - avg. 129.1140, mono. 129.04259
add_M_methionine = 0.0000              # added to M - avg. 131.1961, mono. 131.04048
add_H_histidine = 0.0000               # added to H - avg. 137.1393, mono. 137.05891
add_F_phenylalanine = 0.0000           # added to F - avg. 147.1739, mono. 147.06841
add_R_arginine = 0.0000                # added to R - avg. 156.1857, mono. 156.10111
add_Y_tyrosine = 0.0000                # added to Y - avg. 163.0633, mono. 163.06333
add_W_tryptophan = 0.0000              # added to W - avg. 186.0793, mono. 186.07931
add_B_user_amino_acid = 0.0000         # added to B - avg.   0.0000, mono.   0.00000
add_J_user_amino_acid = 0.0000         # added to J - avg.   0.0000, mono.   0.00000
add_O_user_amino_acid = 0.0000         # added to O - avg.   0.0000, mono    0.00000
add_U_user_amino_acid = 0.0000         # added to U - avg.   0.0000, mono.   0.00000
add_X_user_amino_acid = 0.0000         # added to X - avg.   0.0000, mono.   0.00000
add_Z_user_amino_acid = 0.0000         # added to Z - avg.   0.0000, mono.   0.00000
```

## Setting the Philosopher pipeline configuration file

For the Philosopher analysis we are going to run it using the automated pipeline mode, this mode will automatically run all the necessary steps for us, since we have multiple folders, it would be difficult to run them all manually. For doing so we first need to set the _philosopher.yaml_ configuration file. The configuration file is divided in two sections; the upper part contains a list of all the commands the program is able to automate, the following sections are the individual commands parameter lists. We will set each of the desired commands to _yes_ on the upper part, then we will configure the individual steps. The example below is what we will use for the analysis. You can check on the documentation page the meaning of each parameter and how to adjust them for your analysis.

```
# last updated: 2018-12-21
analytics: false             # reports when a workspace is created for usage estimation (default true)
slackToken:                  # specify the Slack API token
slackChannel:                # specify the channel name
commands:
  workspace: yes             # manage the experiment workspace for the analysis
  database: yes              # target-decoy database formatting
  comet: no                  # peptide spectrum matching with Comet
  msfragger: yes             # peptide spectrum matching with MSFragger
  peptideprophet: yes        # peptide assignment validation
  ptmprophet: no             # PTM site localization
  proteinprophet: no         # protein identification validation
  filter: yes                # statistical filtering, validation and False Discovery Rates assessment
  freequant: yes             # label-free Quantification
  labelquant: yes            # isobaric Labeling-Based Relative Quantification
  report: yes                # multi-level reporting for both narrow-searches and open-searches
  cluster: no                # protein report based on protein clusters
  abacus: yes                # combined analysis of LC-MS/MS results
database:
  add: ''                    # add custom sequences (UniProt FASTA format only)
  annotate: /workspace/6_CPTAC3_Clear_Cell_Renal_Cell_Carcinoma/database/2018-07-16-td-RefSeq.20180629_Human_ucsc_hg38_cpdbnr_mito_264contams.fasta              # process a ready-to-use database
  contam: true               # add common contaminants
  custom: ''                 # use a pre formatted custom database
  enzyme: trypsin            # enzyme for digestion (trypsin, lys_c, lys_n, chymotrypsin) (default "trypsin")
  id: ''                     # UniProt proteome ID
  isoform: false             # add isoform sequences
  prefix: rev_               # decoy prefix to be added (default "rev_")
  reviewed: true             # use only reviewed sequences from Swiss-Prot
comet:
  noindex: true              # skip raw file indexing
  param: ''                  # comet parameter file (default "comet.params.txt")
  raw: mzML                  # format of the spectra file
msfragger:
  path: /workspace/6_CPTAC3_Clear_Cell_Renal_Cell_Carcinoma/bin/MSFragger-20180316.jar                   # path to MSFragger java file
  memmory: 8                 # how much memory in GB to use
  param: /workspace/6_CPTAC3_Clear_Cell_Renal_Cell_Carcinoma/params/fragger.params                  # MSFragger parameter file
  raw: mzML                  # format of the spectra file
peptideprophet:
  extension: pepXML          # pepXML file extension
  accmass: true              # use Accurate Mass model binning
  database: /workspace/6_CPTAC3_Clear_Cell_Renal_Cell_Carcinoma/database/2018-07-16-td-RefSeq.20180629_Human_ucsc_hg38_cpdbnr_mito_264contams.fasta               # path to the database
  decoy: rev_                # semi-supervised mode, protein name prefix to identify Decoy entries
  decoyprobs: true           # compute possible non-zero probabilities for Decoy entries on the last iteration
  enzyme: ''                 # enzyme used in sample (optional)
  exclude: false             # exclude deltaCn*, Mascot*, and Comet* results from results (default Penalize * results)
  expectscore: true          # use expectation value as the only contributor to the f-value for modeling
  forcedistr: false          # bypass quality control checks, report model despite bad modeling
  glyc: false                # enable peptide Glyco motif model
  icat: false                # apply ICAT model (default Autodetect ICAT)
  instrwarn: false           # warn and continue if combined data was generated by different instrument models
  leave: true                # leave alone deltaCn*, Mascot*, and Comet* results from results (default Penalize * results)
  maldi: false               # enable MALDI mode
  masswidth: 5               # model mass width (default 5)
  minpeplen: 7               # minimum peptide length not rejected (default 7)
  minpintt: 2                # minimum number of NTT in a peptide used for positive pI model (default 2)
  minpiprob: 0.9             # minimum probability after first pass of a peptide used for positive pI model (default 0.9)
  minprob: 0.05              # report results with minimum probability (default 0.05)
  minrtntt: 2                # minimum number of NTT in a peptide used for positive RT model (default 2)
  minrtprob: 0.9             # minimum probability after first pass of a peptide used for positive RT model (default 0.9)
  neggamma: false            # use Gamma distribution to model the negative hits
  noicat: false              # do no apply ICAT model (default Autodetect ICAT)
  nomass: false              # disable mass model
  nonmc: false               # disable NMC missed cleavage model
  nonparam: true             # use semi-parametric modeling, must be used in conjunction with --decoy option
  nontt: false               # disable NTT enzymatic termini model
  optimizefval: false        # (SpectraST only) optimize f-value function f(dot,delta) using PCA
  phospho: false             # enable peptide Phospho motif model
  pi: false                  # enable peptide pI model
  ppm: true                  # use PPM mass error instead of Dalton for mass modeling
  rt: false                  # enable peptide RT model
  zero: false                # report results with minimum probability 0
ptmprophet:
  em: 1                      # set EM models to 0 (no EM), 1 (Intensity EM Model Applied) or 2 (Intensity and Matched Peaks EM Models Applied)
  keepold: false             # retain old PTMProphet results in the pepXML file
  verbose: false             # produce Warnings to help troubleshoot potential PTM shuffling or mass difference issues
  mztol: 0.1                 # use specified +/- MS2 mz tolerance on site specific ions
  ppmtol: 1                  # use specified +/- MS1 ppm tolerance on peptides which may have a slight offset depending on search parameters
  minprob: 0                 # use specified minimum probability to evaluate peptides
  massdiffmode: false        # use the Mass Difference and localize
  mods:''                    # specify modifications
proteinprophet:
  accuracy: false            # equivalent to --minprob 0
  allpeps: false             # consider all possible peptides in the database in the confidence model
  confem: false              # use the EM to compute probability given the confidence
  delude: false              # do NOT use peptide degeneracy information when assessing proteins
  excludezeros: false        # exclude zero prob entries
  fpkm: false                # model protein FPKM values
  glyc: false                # highlight peptide N-glycosylation motif
  icat: false                # highlight peptide cysteines
  instances: false           # use Expected Number of Ion Instances to adjust the peptide probabilities prior to NSP adjustment
  iprophet: false            # input is from iProphet
  logprobs: false            # use the log of the probabilities in the Confidence calculations
  maxppmdiff: 1000000        # maximum peptide mass difference in PPM (default 20)
  minprob: 0.05              # peptideProphet probabilty threshold (default 0.05)
  mufactor: 1                # fudge factor to scale MU calculation (default 1)
  nogroupwts: false          # check peptide's Protein weight against the threshold (default: check peptide's Protein Group weight against threshold)
  nonsp: false               # do not use NSP model
  nooccam: false             # non-conservative maximum protein list
  noprotlen: false           # do not report protein length
  normprotlen: false         # normalize NSP using Protein Length
  protmw: false              # get protein mol weights
  softoccam: false           # peptide weights are apportioned equally among proteins within each Protein Group (less conservative protein count estimate)
  unmapped: false            # report results for UNMAPPED proteins
filter:
  pepxml: interact.pep.xml   # overwrites pepXML file NAME (needs to be the same in all directories)
  protxml:                   # overwrites protxml file PATH (needs to be the same in all directories or combined)
  psmFDR: 0.01               # psm FDR level (default 0.01)
  peptideFDR: 0.01           # peptide FDR level (default 0.01)
  ionFDR: 0.01               # peptide ion FDR level (default 0.01)
  proteinFDR: 0.01           # protein FDR level (default 0.01)
  peptideProbability: 0.7    # top peptide probability threshold for the FDR filtering (default 0.7)
  proteinProbability: 0.5    # protein probability threshold for the FDR filtering (not used with the razor algorithm) (default 0.5)
  peptideWeight: 0.9         # threshold for defining peptide uniqueness (default 1)
  tag: rev_                  # decoy tag (default "rev_")
  razor: true                # use razor peptides for protein FDR scoring
  picked: true               # apply the picked FDR algorithm before the protein scoring
  mapMods: false             # map modifications acquired by an open search
  models: true               # print model distribution
  sequential: true           # alternative algorithm that estimates FDR using both filtered PSM and Protein lists
freequant:
  peakTimeWindow: 0.4        # specify the time windows for the peak (minute) (default 0.4)
  retentionTimeWindow: 3     # specify the retention time window for xic (minute) (default 3)
  tolerance: 10              # m/z tolerance in ppm (default 10)
labelquant:
  plex: 10                   # number of channels
  purity: 0.5                # ion purity threshold (default 0.5)
  tolerance: 20              # m/z tolerance in ppm (default 20)
  uniqueOnly: false          # report quantification based on only unique peptides
  bestPSM: true              # select the best PSMs for protein quantification
  removeLow: 0.05            # ignore the lower % of PSMs based on their summed abundances. 0 Means no removal, entry value must be decimal
  minProb: 0.7               # only use PSMs with a minimum probability score
  annotation: annotation.txt             # annotation file with custom names for the TMT channels
report:
  withDecoys: false          # add decoy observations to reports
cluster:
  organismUniProtID: 9606    # UniProt proteome ID
  level: 0.9                 # cluster identity level (default 0.9)
abacus:
  protein: combined.prot.xml # combined protein file
  peptide: ''                # combined peptide file
  tag: rev_                  # decoy tag (default "rev_")
  proteinprobability: 0.9    # minimum protein probability (default 0.9)
  peptideProbability: 0.5    # minimum peptide probability (default 0.5)
  razor: true                # use razor peptides for protein FDR scoring
  picked: true               # apply the picked FDR algorithm before the protein scoring
  uniqueOnly: false          # report TMT quantification based on only unique peptides
  labels: true               # indicates whether the data sets include TMT labels or not
  ```


## Running the pipeline

To start the pipeline we need to run Philosopher using the pipeline command and passing each one of the data sets we wish to process together.

```
$ bin/philosopher pipeline --config params/philosopher.yaml 01CPTAC_CCRCC_W_JHU_20171007 02CPTAC_CCRCC_W_JHU_20171003 03CPTAC_CCRCC_W_JHU_20171022 04CPTAC_CCRCC_W_JHU_20171026 05CPTAC_CCRCC_W_JHU_20171030 06CPTAC_CCRCC_W_JHU_20171120 07CPTAC_CCRCC_W_JHU_20171127 08CPTAC_CCRCC_W_JHU_20171205 09CPTAC_CCRCC_W_JHU_20171215 10CPTAC_CCRCC_W_JHU_20180119 11CPTAC_CCRCC_W_JHU_20180126 12CPTAC_CCRCC_W_JHU_20180202 13CPTAC_CCRCC_W_JHU_20180215 14CPTAC_CCRCC_W_JHU_20180223 15CPTAC_CCRCC_W_JHU_20180315 16CPTAC_CCRCC_W_JHU_20180322 17CPTAC_CCRCC_W_JHU_20180517 18CPTAC_CCRCC_W_JHU_20180521 19CPTAC_CCRCC_W_JHU_20180526 20CPTAC_CCRCC_W_JHU_20180602 21CPTAC_CCRCC_W_JHU_20180621 22CPTAC_CCRCC_W_JHU_20180625 23CPTAC_CCRCC_W_JHU_20180629
```

Each step will be executed consecutively, no other commands are necessary. By the end we will have all individual results for each data set and the combined protein expression matrix containing the all TMT channels converted to the given labels.


## Wrapping up

By the time your analysis is done, you should have different .tsv files in your workspace, those contains the filtered identifications and sequences like PSMs, peptides, ions and proteins.
