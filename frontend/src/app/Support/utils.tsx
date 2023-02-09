import moment from "moment";

// getRangeDates gets the range date time through the current date minus days input
export const getRangeDates = (days: number) => {
    const endDate = new Date();
    const startDate = new Date(new Date().setDate(endDate.getDate() - days))
    startDate.setHours(0,0,0)

    return [moment(startDate), moment(endDate)]
}


// getRangeDateTime gets the range date time through the current date minus hours input
export const getRangeDateTime = (hours: number) => {
    const endDate = new Date();
    const startDate = new Date()
    startDate.setHours(startDate.getHours() - hours, startDate.getMinutes(), startDate.getSeconds())

    return [moment(startDate), moment(endDate)]
}
