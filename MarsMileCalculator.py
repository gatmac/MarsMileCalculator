import sys
import json
import datetime as dt
import logging
import os
import platform

# Append Donors
def appendDonor(donors, name, amount):
    if name in donors:
        donors[name] += amount
    else:
        donors[name] = amount

# Append Miles
def appendMile(miles, donors, donor, ddate):
    while donors[donor] >= mile:
        if len(miles) > 0:
            last = miles[-1][0]
        else:
            last = 0
        miles.append((last + 1, ddate, donor))
        donors[donor] -= mile

# Export Donors
def exportDonors(donors, fname):
    dfile = open(fname, 'wt')
    for d in sorted(donors.keys()):
        dfile.writelines(d + "," + str(donors[d]) + "\n")
    dfile.close()
    logging.info("Exported file " + fname)
 
# Export Donors
def exportMiles(miles, fname):
    dfile = open(fname, 'wt')
    for m in miles:
        dfile.writelines(str(m[0]) + "," + '{:%m/%d/%Y}'.format(m[1]) + "," + m[2] + "\n")
    dfile.close()
    logging.info("Exported file " + fname)

# Main
if __name__ == "__main__":
    calc_log = "MarsMileCalculator.log"
    logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s %(name)-12s %(levelname)-8s %(message)s',
                    datefmt='%m-%d %H:%M',
                    filename=calc_log,
                    filemode='w')

    logging.info("Starting " + sys.argv[0])

    ## Get Platform and Path
    if platform.system() == "Windows":
        ospath = os.path.abspath(os.path.dirname(sys.argv[0])) + "\\"
    else:
        ospath = os.path.abspath(os.path.dirname(sys.argv[0])) + "/"
    logging.info("Using path " + ospath)
    print("Path is", ospath)

    ## Get config data
    if len(sys.argv) == 1:
        cfg_file = ospath + "MarsMileCalculator.json"
        logging.info("Config file not specified, using " + cfg_file)
    else:
        cfg_file = ospath + sys.argv[1]
        logging.info("Using config file " + cfg_file)

    try:
        with open(cfg_file, "r") as file:
            config = json.load(file)
            mile = config['Marsmile']
            fhistory = ospath + config['Donations_In']
            fdonors = ospath + config["Donors_Out"]
            fmiles = ospath + config["Miles_Out"]
    except:
        logging.critical("Unable to read config file.")
        exit(1)

    donors = {}
    mdonors = {} #this is a temporary version of donors used to support Mars Miles.
    miles = []
    dCount = 0

    try:
        history = open(fhistory, 'rt')
        logging.info("Reading donations file " + fhistory + ".")
        header = history.readline().rstrip().split(',')
        #print((header[1]).strip(), header[2].strip())
        for h in history:
            d = h.rstrip().split(',')
            if len(d[1]) > 0:
                appendDonor(donors, d[1], float(d[2]))
                appendDonor(mdonors, d[1], float(d[2]))
                appendMile(miles, mdonors, d[1], dt.datetime.strptime(d[0], '%m/%d/%Y').date())
                dCount += 1
        history.close()
    except:
        logging.critical("Unable to read donations file.")
        exit(1)

    logging.info(str(dCount) + " donations. " + str(len(donors)) + " unique donors. " + str(miles[-1][0]) + " Mars Miles.")
    try:
        exportDonors(donors, fdonors)
    except:
        logging.critical("Unable to write Donors file.")
        exit(1)
    try:
        exportMiles(miles, fmiles)
    except:
        logging.critical("Unable to write Miles file.")
        exit(1)

    print("Output log file " + calc_log + ".")
