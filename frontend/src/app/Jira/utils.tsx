import { Issue } from './ListIssuesTable';

interface Data {
    name: string,
    x: string,
    y: number,
    label: string,
}

export const getIssuesByLabels = (issues: Issue[], labels: string[]) => {
    const filteredLabels = labels.filter(l => l != "all")

    const issuesByLabels = filteredLabels?.map((x) => {
        const count = issues?.filter(y => y["labels"]?.includes(x) == true).length
        return {
            name: "Issues",
            x: x,
            y: count,
            label: "Issues: " + count,
        };
    }).filter(
        (elem, index, arr) => index === arr.findIndex((t) => t.x === elem.x)
    );

    issuesByLabels.sort((a, b) => (a.y < b.y ? -1 : 1));

    return issuesByLabels
}

export const getIssuesByField = (issues: Issue[], label: string, field: string) => {
    const issuesByField = issues?.map((x) => {
        const count = issues?.filter(y => y[field] == x[field] && y["labels"]?.includes(label) == true).length

        return {
            name: label,
            x: x[field],
            y: count,
            label: label + ': ' + count
        };
    }).filter(
        (elem, index, arr) => index === arr.findIndex((t) => t.x === elem.x)
    );

    return issuesByField
}

export const getIssuesByFields = (issues: Issue[], labels: string[], field: string) => {
    const filteredLabels = labels.filter(l => l != "all")
    const issuesByFields: Data[][] = []

    filteredLabels.forEach((label) => {
        const issuesByField = getIssuesByField(issues, label, field)

        issuesByFields.push(issuesByField)

    });

    return issuesByFields
}


interface Legend {
    name: string
}

export const getLegend = (labels) => {
    const filteredLabels = labels.filter(l => l != "all")
    var legend: Legend[] = []

    filteredLabels?.forEach((label) => {
        const l = { name: label }
        legend.push(l)
    })

    return legend
}