# Processing, Filtering and Analyzing a CPTAC 3 cohort


For this example we will see how we processed and analyzed the Clear Cell Renal Carcinoma cohort data from CPTAC 3 using MSFragger, Philosopher and TMT-Integrator. You will learn how to process a large cohort composed by multiple fractionated TMT-labeled data sets. This tutorial will contain as much details as possible for you to reproduce our results on your side. Keep in mind that it is expected to see small differences between some results once that these tools are under constant improvement.

We will need:

* [Philosopher](https://prvst.github.io/philosopher/) (version 20181218 or higher).
* [MSFragger](https://www.nature.com/articles/nmeth.4256) (version 20180316 or higher).
* [Java](http://www.oracle.com/technetwork/java/javase/downloads/jre9-downloads-3848532.html) version 9 (MSFragger requirement).
* The Clear Cell Renal Carcinoma data set from [CPTAC 3](https://cptacdcc.georgetown.edu/cptac/).
* A computer server running GNU/Linux with at least 64GB of RAM.

We ran this example on a Linux Red Hat 7, meaning that the commands you see below will be "Linux compatible", if you are trying to reproduce this on a Windows machine, you will need to adjust the folder separators ('\\' for windows and '/' for Linux).


## Downloading the data set
The CPTAC 3 data is currently available at the NIH / CPTAC Private Data portal, if you are not part of the consortium you will need to sign an agreement that prohibits any publication using this data until the has concluded. We will not need to convert the RAW files because we are using the mzML provided by the consortium.


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
|   ├── TMT-Integrator.jar
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

The annotation file is a simple text file with a map between the TMT channels and the sample labels, this will be useful at the end when we have the final report. Below is an example of the annotation file for the data set #01:

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
database_name = /workspace/6_CPTAC3_Clear_Cell_Renal_Cell_Carcinoma/fasta/RefSeq_20180629/2018-07-16-td-RefSeq.20180629_Human_ucsc_hg38_cpdbnr_mito_264contams.fasta
num_threads = 28                        # 0=poll CPU to set num threads; else specify num threads directly (max 64)

precursor_mass_tolerance = 20.0
precursor_mass_units = 1               # 0=Daltons, 1=ppm

precursor_true_tolerance = 20.00
precursor_true_units = 1
fragment_mass_tolerance = 20.00
fragment_mass_units = 1		       # 0=Daltons, 1=ppm


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
max_variable_mods_per_mod = 3 			# maximum of 5
max_variable_mods_combinations = 50000  	# maximum of 65534, limits number of modified peptides generated from sequence

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
track_zero_topN = 0		       # in addition to topN results, keep track of top results in zero bin
zero_bin_accept_expect = 0	       # boost top zero bin entry to top if it has expect under 0.01 - set to 0 to disable
zero_bin_mult_expect = 1	       # disabled if above passes - multiply expect of zero bin for ordering purposes (does not affect reported expect)
add_topN_complementary = 0

# spectral processing

minimum_peaks = 15                     # required minimum number of peaks in spectrum to search (default 10)
use_topN_peaks = 100
min_fragments_modelling = 3
min_matched_fragments = 4
minimum_ratio = 0.01		       # filter peaks below this fraction of strongest peak
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

```
# last updated: 2018-05-21
analytics: false              # reports when a workspace is created for usage estimation (default true)
slackToken:                   # specify the Slack API token
slackChannel:                 # specify the channel name
commands:
  workspace: yes              # manage the experiment workspace for the analysis
  database: yes               # target-decoy database formatting
  comet: no                  # peptide spectrum matching with Comet
  msfragger: no              # peptide spectrum matching with MSFragger
  peptideprophet: yes         # peptide assignment validation
  ptmprophet: no             # PTM site localization
  proteinprophet: no        # protein identification validation
  filter: yes                 # statistical filtering, validation and False Discovery Rates assessment
  freequant: yes              # label-free Quantification
  labelquant: yes             # isobaric Labeling-Based Relative Quantification
  report: yes                 # multi-level reporting for both narrow-searches and open-searches
  cluster: no                # protein report based on protein clusters
  abacus: yes                 # combined analysis of LC-MS/MS results
database:
  add: ''                    # add custom sequences (UniProt FASTA format only)
  annotate: /storage/cptac/fasta/RefSeq_20180629/2018-07-16-td-RefSeq.20180629_Human_ucsc_hg38_cpdbnr_mito_264contams.fasta              # process a ready-to-use database
  contam: true               # add common contaminants
  custom: ''                 # use a pre formatted custom database
  enzyme: trypsin            # enzyme for digestion (trypsin, lys_c, lys_n, chymotrypsin) (default "trypsin")
  id: ''                     # UniProt proteome ID
  isoform: false             # add isoform sequences
  prefix: rev_               # decoy prefix to be added (default "rev_")
  reviewed: true             # use only reviwed sequences from Swiss-Prot
comet:
  noindex: true              # skip raw file indexing
  param: ''                  # comet parameter file (default "comet.params.txt")
  raw: mzML                  # format of the spectra file
msfragger:
  path: ''                   # path to MSFragger java file
  memmory: 8                 # how much memmory in GB to use
  param: ''                  # MSFragger parameter file
  raw: mzML                  # format of the spectra file
peptideprophet:
  extension: pepXML          # pepXML file extension
  accmass: true             # use Accurate Mass model binning
  database: /storage/cptac/fasta/RefSeq_20180629/2018-07-16-td-RefSeq.20180629_Human_ucsc_hg38_cpdbnr_mito_264contams.fasta               # path to the database
  decoy: rev_                # semi-supervised mode, protein name prefix to identify Decoy entries
  decoyprobs: true          # compute possible non-zero probabilities for Decoy entries on the last iteration
  enzyme: ''                 # enzyme used in sample (optional)
  exclude: false             # exclude deltaCn*, Mascot*, and Comet* results from results (default Penalize * results)
  expectscore: true         # use expectation value as the only contributor to the f-value for modeling
  forcedistr: false          # bypass quality control checks, report model despite bad modelling
  glyc: false                # enable peptide Glyco motif model
  icat: false                # apply ICAT model (default Autodetect ICAT)
  instrwarn: false           # warn and continue if combined data was generated by different instrument models
  leave: true               # leave alone deltaCn*, Mascot*, and Comet* results from results (default Penalize * results)
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
  nonparam: true            # use semi-parametric modeling, must be used in conjunction with --decoy option
  nontt: false               # disable NTT enzymatic termini model
  optimizefval: false        # (SpectraST only) optimize f-value function f(dot,delta) using PCA
  phospho: false             # enable peptide Phospho motif model
  pi: false                  # enable peptide pI model
  ppm: true                 # use PPM mass error instead of Daltons for mass modeling
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
  mods: STY:79.9663,M:15.9949,S:229.16293,n:229.16293,n:42.01 # specify modifications
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
  maxppmdiff: 1000000             # maximum peptide mass difference in PPM (default 20)
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
  pepxml: interact.pep.xml                   # overwrites pepXML file NAME (needs to be the same in all directories)
  protxml: /storage/cptac/6_CPTAC3_Clear_Cell_Renal_Cell_Carcinoma/phospho/combined.prot.xml                  # overwrites protxml file PATH (needs to be the same in all directories or combined)
  psmFDR: 0.01               # psm FDR level (default 0.01)
  peptideFDR: 0.01           # peptide FDR level (default 0.01)
  ionFDR: 0.01               # peptide ion FDR level (default 0.01)
  proteinFDR: 0.01           # protein FDR level (default 0.01)
  peptideProbability: 0.7    # top peptide probability treshold for the FDR filtering (default 0.7)
  proteinProbability: 0.5    # protein probability treshold for the FDR filtering (not used with the razor algorithm) (default 0.5)
  peptideWeight: 0.9           # threshold for defining peptide uniqueness (default 1)
  tag: rev_                  # decoy tag (default "rev_")
  razor: true               # use razor peptides for protein FDR scoring
  picked: true              # apply the picked FDR algorithm before the protein scoring
  mapMods: false             # map modifications aquired by an open search
  models: true              # print model distribution
  sequential: true          # alternative algorithm that estimates FDR using both filtered PSM and Protein lists
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
  minProb: 0.7               # only use PSMs with a minimun probability score
  annotation: ''             # annotation file with custom names for the TMT channels
report:
  withDecoys: false          # add decoy observations to reports
cluster:
  organismUniProtID: 9606      # UniProt proteome ID
  level: 0.9                 # cluster identity level (default 0.9)
abacus:
  protein: combined.prot.xml # combined protein file
  peptide: combined.pep.xml  # combined peptide file
  tag: rev_                  # decoy tag (default "rev_")
  proteinprobability: 0.9    # minimun protein probability (default 0.9)
  peptideProbability: 0.5    # minimun peptide probability (default 0.5)
  razor: true               # use razor peptides for protein FDR scoring
  picked: true              # apply the picked FDR algorithm before the protein scoring
  uniqueOnly: false          # report TMT quantification based on only unique peptides
  labels: true              # indicates wether the data sets incluse TMT labels or not
  ```

## Creating a Workspace
Now that we have all files in place, we need to transform our directory into a Philosopher [workspace](workspace.md), using the command terminal, type:

```
$ bin/philosopher workspace --init

INFO[12:01:16] Executing Workspace 20180319
INFO[12:01:16] Creating workspace
INFO[12:01:17] Done
```

Every time you see an output on the screen it will start with the name of the command and the version of Philosopher, in my case it tells me i just ran the `workspace` command using Philosopher version 20180319.
The `workspace` command creates a local directory called `.meta` (invisible in GNU/Linux), this directory will contain all processed data for the analysis in a custom binary format for fast access.


## Fetching a human protein database
There are different ways of doing this, but we suggest you to use the Philosopher [database](database.md) command:

```
$ bin/philosopher database --id UP000005640 --reviewed --contam

INFO[12:09:32] Executing Database 20180319
INFO[12:09:32] Fetching database
INFO[12:09:48] Processing decoys
INFO[12:09:48] Creating file
INFO[12:09:59] Done
```

The `database` command connects to [UniProt](https://www.uniprot.org) and fetches a bleeding edge snapshot of an organism proteome and it uses the downloaded file to create a target-decoy database. The parameter that tells Philosopher whats the organism you want is the `--id`, that code you see there is called [Proteome Identifier](https://www.uniprot.org/help/proteome_id), to find it you have to go to the [Proteomes](http://www.uniprot.org/proteomes/) page on [UniProt](https://www.uniprot.org) and query your organism ([UP000005640](http://www.uniprot.org/proteomes/UP000005640) is for _Homo sapiens_).

In our example we are fetching the _Homo sapiens_ proteome containing only reviewed entries (`--reviewed`) and we are adding common contaminants (`--contam`). Our workspace now contains a protein FASTA database file.

```
.
├── 2018-03-21-td-UP000005640.fas
├── b1906_293T_proteinID_01A_QE3_122212.mzML
├── b1922_293T_proteinID_02A_QE3_122212.mzML
├── b1923_293T_proteinID_03A_QE3_122212.mzML
├── b1924_293T_proteinID_04A_QE3_122212.mzML
├── b1925_293T_proteinID_05A_QE3_122212.mzML
├── b1926_293T_proteinID_06A_QE3_122212.mzML
├── b1927_293T_proteinID_07A_QE3_122212.mzML
├── b1928_293T_proteinID_08A_QE3_122212.mzML
├── b1929_293T_proteinID_09A_QE3_122212.mzML
├── b1930_293T_proteinID_10A_QE3_122212.mzML
├── b1931_293T_proteinID_11A_QE3_122212.mzML
├── b1932_293T_proteinID_12A_QE3_122212.mzML
├── b1937_293T_proteinID_01B_QE3_122212.mzML
├── b1938_293T_proteinID_02B_QE3_122212.mzML
├── b1939_293T_proteinID_03B_QE3_122212.mzML
├── b1940_293T_proteinID_04B_QE3_122212.mzML
├── b1941_293T_proteinID_05B_QE3_122212.mzML
├── b1942_293T_proteinID_06B_QE3_122212.mzML
├── b1943_293T_proteinID_07B_QE3_122212.mzML
├── b1944_293T_proteinID_08B_QE3_122212.mzML
├── b1945_293T_proteinID_09B_QE3_122212.mzML
├── b1946_293T_proteinID_10B_QE3_122212.mzML
├── b1947_293T_proteinID_11B_QE3_122212.mzML
├── b1948_293T_proteinID_12B_QE3_122212.mzML
├── bin
│   ├── MSFragger-20180316.jar
│   └── philosopher
└── params
    └── fragger.params

2 directories, 28 files
```

## Doing an open database search using MSFragger
The first thing to do before running the database search is to properly set the parameter file, this is the one I'm using for this example:

```
database_name = /home/felipevl/tutorial/2018-03-21-td-UP000005640.fas
num_threads = 8                        # 0=poll CPU to set num threads; else specify num threads directly (max 64)

precursor_mass_lower = -500.00
precursor_mass_upper = 500.00
precursor_mass_units = 0               # 0=Daltons, 1=ppm
precursor_true_tolerance = 20.00
precursor_true_units = 1               # 0=Daltons, 1=ppm
fragment_mass_tolerance = 20.00
fragment_mass_units = 1		        # 0=Daltons, 1=ppm

isotope_error = 0                      # 0=off, 0/1/2 (standard C13 error)
mass_offsets = 0					   # 0=standard

search_enzyme_name = Trypsin
search_enzyme_cutafter = KR
search_enzyme_butnotafter = P

num_enzyme_termini = 2                 # 2 for enzymatic, 1 for semi-enzymatic, 0 for nonspecific digestion
allowed_missed_cleavage = 1            # maximum value is 5

clip_nTerm_M = 1

#maximum of 7 mods - amino acid codes, * for any amino acid, ^ for termini, [ and ] specifies protein termini, n and c specifies peptide termini
variable_mod_01 = 15.9949 M
variable_mod_02 = 42.0106 [^
#variable_mod_03 = 79.96633 STY
#variable_mod_03 = -17.0265 nQnC
#variable_mod_04 = -18.0106 nE

allow_multiple_variable_mods_on_residue = 1  	   # static mods are not considered
max_variable_mods_per_mod = 3 			          # maximum of 5
max_variable_mods_combinations = 5000  		     # maximum of 65534, limits number of modified peptides generated from sequence

output_file_extension = pepXML	             # pepXML
output_format = pepXML				         # pepXML or tsv
output_report_topN = 1
output_max_expect = 50

precursor_charge = 0 0                 # precursor charge range to analyze; does not override any existing charge; 0 as 1st entry ignores parameter
override_charge = 0                    # 0=no, 1=yes to override existing precursor charge states with precursor_charge parameter

digest_min_length = 7
digest_max_length = 50
digest_mass_range = 500.0 5000.0       # MH+ peptide mass range to analyze
max_fragment_charge = 2                # set maximum fragment charge state to analyze (allowed max 5)

#open search parameters
track_zero_topN = 0		           # in addition to topN results, keep track of top results in zero bin
zero_bin_accept_expect = 0	        # boost top zero bin entry to top if it has expect under 0.01 - set to 0 to disable
zero_bin_mult_expect = 1	          # disabled if above passes - multiply expect of zero bin for ordering purposes (does not affect reported expect)
add_topN_complementary = 0

# spectral processing

minimum_peaks = 15                   # required minimum number of peaks in spectrum to search (default 10)
use_topN_peaks = 100
min_fragments_modelling = 3
min_matched_fragments = 6
minimum_ratio = 0.01		         # filter peaks below this fraction of strongest peak
clear_mz_range = 0.0 0.0             # for iTRAQ/TMT type data; will clear out all peaks in the specified m/z range

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
add_K_lysine = 0.0000                  # added to K - avg. 128.1723, mono. 128.09496
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

There are several details here that you have to pay attention before running the analysis. The first thing is to use the correct _full_ path to your recently downloaded FASTA database, the example above contains the path to mine on mine machine, so you need to replace that. Make sure that you have the `precursor_mass_lower` and `precursor_mass_upper` set to -500 and 500 respectively, and that the units are set to Daltons. I'm not going to explain all the details on the search parameter because it's not my goal right now, you can find all the necessary information on the software [manual](http://inventions.umich.edu/technologies/7143_msfragger-ultrafast-and-comprehensive-identification-of-peptides-from-tandem-mass-spectra).

Now that everything is set in place we can run MSFragger:

```
$ java -jar -Xmx20G bin/MSFragger-20180316.jar params/fragger.params *.mzML

MSFragger version MSFragger-20180316

MSFTBX version 1.8.2
(c) University of Michigan

System OS: Linux, Architecture: amd64
Java Info: 1.8.0_151, OpenJDK 64-Bit Server VM, Oracle Corporation
JVM started with 18204MB memory

Sequence database filtered and tagged in 45ms
Digestion completed in 269ms
Merged digestion results in 80ms
Sorting digested sequences...
        of length 7: 239424
        of length 8: 217125
        of length 9: 213064
        of length 10: 189236
        of length 11: 173845
        of length 12: 159615
        of length 13: 150149
        of length 14: 134445
        of length 15: 126520
        of length 16: 114174
        of length 17: 104548
        of length 18: 97248
        of length 19: 91045
        of length 20: 82447
        of length 21: 77214
        of length 22: 70461
        of length 23: 65701
        of length 24: 61205
        of length 25: 56870
        of length 26: 52014
        of length 27: 49177
        of length 28: 45802
        of length 29: 42139
        of length 30: 39333
        of length 31: 36062
        of length 32: 33895
        of length 33: 30978
        of length 34: 28733
        of length 35: 26073
        of length 36: 24913
        of length 37: 22199
        of length 38: 20595
        of length 39: 18310
        of length 40: 17647
        of length 41: 16516
        of length 42: 14535
        of length 43: 12755
        of length 44: 10316
        of length 45: 8123
        of length 46: 5382
        of length 47: 3500
        of length 48: 2163
        of length 49: 1159
        of length 50: 707
        DONE
Removing duplicates and compacting...
Reduced to 2836180  peptides in 3424ms
Generating modified peptides...DONE in 598ms
Generated 4354434 modified peptides
Merging peptide pools from threads... DONE in 69ms
Sorting modified peptides by mass...DONE in 358ms
Peptide index written in 137ms
Selected fragment tolerance 0.02 Da and maximum fragment slice size of 13704.53MB
323824040 fragments to be searched in 1 slices (2.41GB total)
Operating on slice 1 of 1: 10656ms
        b1906_293T_proteinID_01A_QE3_122212.mzML 9726ms [progress: 41820/41820 (100.00%) - 116.95 spectra/s] - completed 285675ms
        b1922_293T_proteinID_02A_QE3_122212.mzML 8264ms [progress: 46856/46856 (100.00%) - 126.80 spectra/s] - completed 325529ms
        b1923_293T_proteinID_03A_QE3_122212.mzML 9141ms [progress: 42672/42672 (100.00%) - 34.82 spectra/s] - completed 288781ms
        b1924_293T_proteinID_04A_QE3_122212.mzML 7394ms [progress: 48033/48033 (100.00%) - 76.69 spectra/s] - completed 319507ms
        b1925_293T_proteinID_05A_QE3_122212.mzML 7281ms [progress: 48788/48788 (100.00%) - 108.59 spectra/s] - completed 307355ms
        b1926_293T_proteinID_06A_QE3_122212.mzML 7497ms [progress: 42828/42828 (100.00%) - 163.70 spectra/s] - completed 274925ms
        b1927_293T_proteinID_07A_QE3_122212.mzML 7372ms [progress: 52252/52252 (100.00%) - 114.62 spectra/s] - completed 339057ms
        b1928_293T_proteinID_08A_QE3_122212.mzML 7272ms [progress: 44054/44054 (100.00%) - 38.33 spectra/s] - completed 281215ms
        b1929_293T_proteinID_09A_QE3_122212.mzML 7730ms [progress: 58817/58817 (100.00%) - 67.21 spectra/s] - completed 388199ms
        b1930_293T_proteinID_10A_QE3_122212.mzML 8070ms [progress: 45131/45131 (100.00%) - 42.18 spectra/s] - completed 281539ms
        b1931_293T_proteinID_11A_QE3_122212.mzML 7790ms [progress: 55352/55352 (100.00%) - 74.56 spectra/s] - completed 357178ms
        b1932_293T_proteinID_12A_QE3_122212.mzML 7310ms [progress: 41313/41313 (100.00%) - 147.78 spectra/s] - completed 244440ms
        b1937_293T_proteinID_01B_QE3_122212.mzML 7665ms [progress: 46816/46816 (100.00%) - 121.83 spectra/s] - completed 288247ms
        b1938_293T_proteinID_02B_QE3_122212.mzML 7467ms [progress: 45642/45642 (100.00%) - 41.13 spectra/s] - completed 277271ms
        b1939_293T_proteinID_03B_QE3_122212.mzML 7203ms [progress: 51281/51281 (100.00%) - 31.38 spectra/s] - completed 321343ms
        b1940_293T_proteinID_04B_QE3_122212.mzML 7382ms [progress: 49014/49014 (100.00%) - 100.60 spectra/s] - completed 293743ms
        b1941_293T_proteinID_05B_QE3_122212.mzML 7548ms [progress: 44128/44128 (100.00%) - 36.89 spectra/s] - completed 271936ms
        b1942_293T_proteinID_06B_QE3_122212.mzML 7709ms [progress: 54072/54072 (100.00%) - 35.00 spectra/s] - completed 341091ms
        b1943_293T_proteinID_07B_QE3_122212.mzML 9119ms [progress: 39929/39929 (100.00%) - 23.33 spectra/s] - completed 250896ms
        b1944_293T_proteinID_08B_QE3_122212.mzML 8377ms [progress: 58487/58487 (100.00%) - 88.85 spectra/s] - completed 382977ms
        b1945_293T_proteinID_09B_QE3_122212.mzML 7303ms [progress: 17752/17752 (100.00%) - 202.16 spectra/s] - completed 100209ms
        b1946_293T_proteinID_10B_QE3_122212.mzML 9737ms [progress: 57882/57882 (100.00%) - 59.12 spectra/s] - completed 437576ms
        b1947_293T_proteinID_11B_QE3_122212.mzML 8092ms [progress: 33042/33042 (100.00%) - 164.68 spectra/s] - completed 215187ms
        b1948_293T_proteinID_12B_QE3_122212.mzML 9622ms [progress: 55188/55188 (100.00%) - 43.96 spectra/s] - completed 412216ms
```

After completion, we now have one `pepXML` file for each `mzML`:

```
.
├── 2018-03-21-td-UP000005640.fas
├── 2018-03-21-td-UP000005640.fas.1.pepindex
├── b1906_293T_proteinID_01A_QE3_122212.mzML
├── b1906_293T_proteinID_01A_QE3_122212.pepXML
├── b1922_293T_proteinID_02A_QE3_122212.mzML
├── b1922_293T_proteinID_02A_QE3_122212.pepXML
├── b1923_293T_proteinID_03A_QE3_122212.mzML
├── b1923_293T_proteinID_03A_QE3_122212.pepXML
├── b1924_293T_proteinID_04A_QE3_122212.mzML
├── b1924_293T_proteinID_04A_QE3_122212.pepXML
├── b1925_293T_proteinID_05A_QE3_122212.mzML
├── b1925_293T_proteinID_05A_QE3_122212.pepXML
├── b1926_293T_proteinID_06A_QE3_122212.mzML
├── b1926_293T_proteinID_06A_QE3_122212.pepXML
├── b1927_293T_proteinID_07A_QE3_122212.mzML
├── b1927_293T_proteinID_07A_QE3_122212.pepXML
├── b1928_293T_proteinID_08A_QE3_122212.mzML
├── b1928_293T_proteinID_08A_QE3_122212.pepXML
├── b1929_293T_proteinID_09A_QE3_122212.mzML
├── b1929_293T_proteinID_09A_QE3_122212.pepXML
├── b1930_293T_proteinID_10A_QE3_122212.mzML
├── b1930_293T_proteinID_10A_QE3_122212.pepXML
├── b1931_293T_proteinID_11A_QE3_122212.mzML
├── b1931_293T_proteinID_11A_QE3_122212.pepXML
├── b1932_293T_proteinID_12A_QE3_122212.mzML
├── b1932_293T_proteinID_12A_QE3_122212.pepXML
├── b1937_293T_proteinID_01B_QE3_122212.mzML
├── b1937_293T_proteinID_01B_QE3_122212.pepXML
├── b1938_293T_proteinID_02B_QE3_122212.mzML
├── b1938_293T_proteinID_02B_QE3_122212.pepXML
├── b1939_293T_proteinID_03B_QE3_122212.mzML
├── b1939_293T_proteinID_03B_QE3_122212.pepXML
├── b1940_293T_proteinID_04B_QE3_122212.mzML
├── b1940_293T_proteinID_04B_QE3_122212.pepXML
├── b1941_293T_proteinID_05B_QE3_122212.mzML
├── b1941_293T_proteinID_05B_QE3_122212.pepXML
├── b1942_293T_proteinID_06B_QE3_122212.mzML
├── b1942_293T_proteinID_06B_QE3_122212.pepXML
├── b1943_293T_proteinID_07B_QE3_122212.mzML
├── b1943_293T_proteinID_07B_QE3_122212.pepXML
├── b1944_293T_proteinID_08B_QE3_122212.mzML
├── b1944_293T_proteinID_08B_QE3_122212.pepXML
├── b1945_293T_proteinID_09B_QE3_122212.mzML
├── b1945_293T_proteinID_09B_QE3_122212.pepXML
├── b1946_293T_proteinID_10B_QE3_122212.mzML
├── b1946_293T_proteinID_10B_QE3_122212.pepXML
├── b1947_293T_proteinID_11B_QE3_122212.mzML
├── b1947_293T_proteinID_11B_QE3_122212.pepXML
├── b1948_293T_proteinID_12B_QE3_122212.mzML
├── b1948_293T_proteinID_12B_QE3_122212.pepXML
├── bin
│   ├── MSFragger-20180316.jar
│   └── philosopher
└── params
    └── fragger.params
```


## Validating the peptide identifications
Having finished with the database search, now we need to run the [PeptidePrphet](peptideprophet.md) command to validate and assign confidence scores to the PSMs, this will allow us to `filter` the PSMs using a false discovery rate.

```
$ bin/philosopher peptideprophet --accmass --clevel -2 --combine --database 2018-03-21-td-UP000005640.fas --decoy rev_ --decoyprobs --expectscore --nonparam --masswidth 1000.0 *.pepXML

(...)

INFO: Processing standard MixtureModel ...

 PeptideProphet  (TPP v5.0.1 Post-Typhoon dev, Build 201705191541-exported (Linux-x86_64)) AKeller@ISB
 read in 0 1+, 597664 2+, 347529 3+, 80435 4+, 8775 5+, 2169 6+, and 86 7+ spectra.
Initializing statistical models ...
Found 136993 Decoys, and 899665 Non-Decoys
Iterations: .........10.........20......
WARNING: Mixture model quality test failed for charge (1+).
model complete after 27 iterations
INFO[15:43:40] Done
```

PeptideProphet output contains a lot of log information that we don't need to focus right now. Pay attention to the number of peptides found within the charge state distribution and the number of target and decoys matches, if they are not the same I'm showing up here they should be at least very close to that. The output from this command is a new filed called `interact.pep.xml`, it contains the combined validation from all your `pepXML` files and that's what we are going to use from now on.


## Inferring proteins from peptide matches
For the protein inference and validation we will use the [ProteinProphet](proteinprophet.md) command as follows:

```
$ bin/philosopher proteinprophet --maxppmdiff 100000 interact.pep.xml

INFO[15:48:54] Executing ProteinProphet 20180319
ProteinProphet (C++) by Insilicos LLC and LabKey Software, after the original Perl by A. Keller (TPP v5.0.1 Post-Typhoon dev, Build 201705221614-exported (Linux-x86_64))
 (no FPKM) (using degen pep info)
Reading in /home/felipevl/tutorial/interact.pep.xml...
...read in 0 1+, 398519 2+, 213862 3+, 39588 4+, 3827 5+, 326 6+, 23 7+ spectra with min prob 0.05

(...)

Finished
INFO[15:59:22] Done
```

In a similar way to PeptideProphet, the output from ProteinProphet is another file called `interact.prot.xml`.


## Calculating false discovery rates (FDR)
Philosopher allows us to calculate multi-level false discovery rates through the [Filter](filter.md) command. Multi-level means that our data will be interpreted in different ways; PSMs, peptides, peptide ions and proteins, and for each layer, an FDR score will be calculated simultaneously.

The `filter` command is very versatile, it allows you to explore your data using different approaches and algorithms, you can read more about all the options available in the [documentation](filter.md). Here we are going to focus on getting similar results to the ones published before.

```
$ bin/philosopher filter --pepxml interact.pep.xml --protxml interact.prot.xml --sequential --razor --mapmods --models

INFO[16:09:16] Executing Filter 20180319
INFO[16:09:16] Processing peptide identification files
INFO[16:10:16] 1+ Charge profile                             decoy=0 target=0
INFO[16:10:16] 2+ Charge profile                             decoy=11959 target=420104
INFO[16:10:16] 3+ Charge profile                             decoy=6055 target=226990
INFO[16:10:17] 4+ Charge profile                             decoy=1669 target=43507
INFO[16:10:17] 5+ Charge profile                             decoy=201 target=4094
INFO[16:10:17] 6+ Charge profile                             decoy=45 target=347
INFO[16:10:21] Database search results                       ions=217892 peptides=156182 psms=714999
INFO[16:10:25] Converged to 1.00 % FDR with 647474 PSMs      decoy=6535 threshold=0.3543 total=654009
INFO[16:10:29] Converged to 1.00 % FDR with 116627 Peptides  decoy=1176 threshold=0.7649 total=117803
INFO[16:10:32] Converged to 1.00 % FDR with 171494 Ions      decoy=1731 threshold=0.7153 total=173225
INFO[16:10:40] Protein inference results                     decoy=5677 target=12890
INFO[16:10:41] Converged to 1.01 % FDR with 9439 Proteins    decoy=95 threshold=0.9873 total=9534
INFO[16:10:48] Applying sequential FDR estimation            ions=183803 peptides=127261 psms=646288
INFO[16:10:52] Converged to 0.25 % FDR with 644654 PSMs      decoy=1634 threshold=0.3543 total=646288
INFO[16:10:55] Converged to 0.52 % FDR with 126602 Peptides  decoy=659 threshold=0.3543 total=127261
INFO[16:10:57] Converged to 0.38 % FDR with 183106 Ions      decoy=697 threshold=0.3543 total=183803
INFO[16:10:58] Post processing identifications
INFO[16:11:05] Mapping modifications
INFO[16:11:44] Processing modifications
INFO[16:12:23] Processing Protein Inference
INFO[16:12:45] Correcting PSM to Protein mappings
INFO[16:21:08] Mapping Ion status to PSMs
INFO[16:21:11] Updating razor PSM assingment to Proteins
INFO[16:21:13] Calculating Spectral Counts
INFO[16:35:52] Done
```

A key parameter for open search results is `--mapmods`, this option will organize and process the identifications based on the calculated delta masses. Philosopher will try to match those to entries from [UniMod](http://www.unimod.org/).

According to our results we have:

* 644654 PSMs
* 126602 Peptides
* 183106 Ions
* 9439 proteins

## Reporting
Lastly, in order to inspect our results we need to run the [Report](report.md) command.

```
$ bin/philosopher report

INFO[16:38:18] Executing Report 20180319
INFO[16:39:08] Creating Protein FASTA report
INFO[16:39:08] Creating Protein report
INFO[16:39:09] Creating PSM report
INFO[16:39:18] Creating peptide report
INFO[16:39:18] Creating peptide Ion report
INFO[16:39:20] Creating modification reports
INFO[16:39:20] Plotting mass distribution
INFO[16:39:20] Done
```

This will create different reports based on the processed layers plus reports regarding the modification profile. This is what our workspace looks like now:

```

├── 2018-03-21-td-UP000005640.fas
├── 2018-03-21-td-UP000005640.fas.1.pepindex
├── b1906_293T_proteinID_01A_QE3_122212.mzML
├── b1906_293T_proteinID_01A_QE3_122212.pepXML
├── b1922_293T_proteinID_02A_QE3_122212.mzML
├── b1922_293T_proteinID_02A_QE3_122212.pepXML
├── b1923_293T_proteinID_03A_QE3_122212.mzML
├── b1923_293T_proteinID_03A_QE3_122212.pepXML
├── b1924_293T_proteinID_04A_QE3_122212.mzML
├── b1924_293T_proteinID_04A_QE3_122212.pepXML
├── b1925_293T_proteinID_05A_QE3_122212.mzML
├── b1925_293T_proteinID_05A_QE3_122212.pepXML
├── b1926_293T_proteinID_06A_QE3_122212.mzML
├── b1926_293T_proteinID_06A_QE3_122212.pepXML
├── b1927_293T_proteinID_07A_QE3_122212.mzML
├── b1927_293T_proteinID_07A_QE3_122212.pepXML
├── b1928_293T_proteinID_08A_QE3_122212.mzML
├── b1928_293T_proteinID_08A_QE3_122212.pepXML
├── b1929_293T_proteinID_09A_QE3_122212.mzML
├── b1929_293T_proteinID_09A_QE3_122212.pepXML
├── b1930_293T_proteinID_10A_QE3_122212.mzML
├── b1930_293T_proteinID_10A_QE3_122212.pepXML
├── b1931_293T_proteinID_11A_QE3_122212.mzML
├── b1931_293T_proteinID_11A_QE3_122212.pepXML
├── b1932_293T_proteinID_12A_QE3_122212.mzML
├── b1932_293T_proteinID_12A_QE3_122212.pepXML
├── b1937_293T_proteinID_01B_QE3_122212.mzML
├── b1937_293T_proteinID_01B_QE3_122212.pepXML
├── b1938_293T_proteinID_02B_QE3_122212.mzML
├── b1938_293T_proteinID_02B_QE3_122212.pepXML
├── b1939_293T_proteinID_03B_QE3_122212.mzML
├── b1939_293T_proteinID_03B_QE3_122212.pepXML
├── b1940_293T_proteinID_04B_QE3_122212.mzML
├── b1940_293T_proteinID_04B_QE3_122212.pepXML
├── b1941_293T_proteinID_05B_QE3_122212.mzML
├── b1941_293T_proteinID_05B_QE3_122212.pepXML
├── b1942_293T_proteinID_06B_QE3_122212.mzML
├── b1942_293T_proteinID_06B_QE3_122212.pepXML
├── b1943_293T_proteinID_07B_QE3_122212.mzML
├── b1943_293T_proteinID_07B_QE3_122212.pepXML
├── b1944_293T_proteinID_08B_QE3_122212.mzML
├── b1944_293T_proteinID_08B_QE3_122212.pepXML
├── b1945_293T_proteinID_09B_QE3_122212.mzML
├── b1945_293T_proteinID_09B_QE3_122212.pepXML
├── b1946_293T_proteinID_10B_QE3_122212.mzML
├── b1946_293T_proteinID_10B_QE3_122212.pepXML
├── b1947_293T_proteinID_11B_QE3_122212.mzML
├── b1947_293T_proteinID_11B_QE3_122212.pepXML
├── b1948_293T_proteinID_12B_QE3_122212.mzML
├── b1948_293T_proteinID_12B_QE3_122212.pepXML
├── delta-mass.html
├── interact.pep.xml
├── interact.prot.xml
├── ion.tsv
├── modifications.tsv
├── peptide.tsv
├── proteins.fas
├── psm.tsv
├── report.tsv
├── bin
│   ├── MSFragger-20180316.jar
│   └── philosopher
├── params
    └── fragger.params

2 directories, 62 files
```

## Wrapping up

By the time your analysis is done, you should have different .tsv files in your workspace, those contains the filtered identifications and sequences like PSMs, peptides, ions and proteins.
