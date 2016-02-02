# This is a program that makes it easy to set WFH updates
Really, at this point, it's just a script and it is a serious misuse of the Go language

#Set up
1. Set up Go. Maybe use brew or something idgaf
2. [Get yourself some API secrets from google](https://console.developers.google.com/flows/enableapi?apiid=calendar) (dude, they give them to you FOR FREE)
3. Set your PRODUCT_WFH_CALENDAR_ID environment variable in your ~/.bash_profile to the thing that @mkuipers tells you
3. Set your WFH_DISPLAY_CALENDAR_NAME environment variable in your ~/.bash_profile to "#{first_name} {last_initial}" (e.g. Max K)

#Usage
This is designed to be used over the weekend prior to week in which you'd be setting your schedule or that very week.  It only knows about days that are the current date + 6 days.  

I don't think you'll need to build it yourself because I don't really know much about how Go works but if you do, then go ahead and do that. Otherwise you can just use the executable in this directory.

Then do like

```
./wfh tue
```
Or whatever 3-letter abbreviation of a day.  
