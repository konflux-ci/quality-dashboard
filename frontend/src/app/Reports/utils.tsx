import { yyyyMMddFormat } from "@patternfly/react-core";
import moment from 'moment'

// getRangeDates gets the range date time through the current date minus days input
export const getRangeDates = (days: number) => {
    const endDate = new Date();
    const startDate = new Date(new Date().setDate(endDate.getDate() - days))
    startDate.setHours(0, 0, 0)
    endDate.setSeconds(0)

    return [startDate, endDate]
}

// getRangeDateTime gets the range date time through the current date minus hours input
export const getRangeDateTime = (hours: number) => {
    const endDate = new Date();
    const startDate = new Date()
    startDate.setHours(startDate.getHours() - hours, startDate.getMinutes(), startDate.getSeconds())

    return [startDate, endDate]
}

// getTimes gets the time from a date with format HH:mm:ss
export const getTime = (date: Date) => {
    return date.getHours().toString().padStart(2, '0') + 
    ":" + date.getMinutes().toString().padStart(2, '0') + 
        ":" + date.getSeconds().toString().padStart(2, '0')
}

// getDateFormat gets date format based on user's locale
export const getDateFormat = () => {
    return moment.localeData(navigator.language).longDateFormat('L');
}

// parseCustomizedDate returns date parsed based on users locale
export const parseCustomizedDate = (dateString: string) => {
    return moment(dateString, getDateFormat()).toDate();
}

// customizedFormatDate returns date formated based on users locale
export const customizedFormatDate = (date: Date) => {
    return new Intl.DateTimeFormat(navigator.language).format(date);
}

// customizedFormatDateTime returns date and time formated based on users locale
export const customizedFormatDateTime = (date: Date) => {
    return customizedFormatDate(date) + " " + getTime(date);
}

export const validateCustomizedDate = (dateString: string) => {
    var dateFormat = getDateFormat()
    return moment(moment(dateString).format(dateFormat),dateFormat,true).isValid() ? '' : 'Invalid time format';
}

// formatDate return a date formated with yyyyMMdd format
export const formatDate = (date: Date) => {
    return yyyyMMddFormat(date);
}

// formatDate return a date formated with yyyyMMdd format and time
export const formatDateTime = (date: Date) => {
    return yyyyMMddFormat(date) + " " + getTime(date);
}

// rangeInfo contains all the necessary info to define quick start ranges
export type rangeInfo = {
    type: string,
    days: number,
    hours: number,
}

// ranges contains all the search quick ranges
export const ranges = [
    { type: "Last 12 hours", days: 0, hours: 12 },
    { type: "Last day", days: 1, hours: 0 },
    { type: "Last 2 days", days: 2, hours: 0 },
    { type: "Last 1 week", days: 7, hours: 0 },
    { type: "Last 2 weeks", days: 15, hours: 0 },
    { type: "Last month", days: 30, hours: 0 },
]
