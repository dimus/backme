# ``backme`` --- A backup files organizer

Quite often big files (like database dumps) are backed up on a
regular basis. Such backups become a disk space hogs and when time comes to
restoring data it is cumbersome to find the right file for the job. This
script checks files of a certain pattern and sorts them by date into daily
monthly, yearly directories, progressively deleting more and more of aged
files. For example it keeps all files from the last 24 hours, but keeps only
one file per day in the last 30 days, only one file per month for the first
3 years, and after that only 1 file per year. As a result backup takes
significantly less space and it is easier to find a file from a specific
period of time.