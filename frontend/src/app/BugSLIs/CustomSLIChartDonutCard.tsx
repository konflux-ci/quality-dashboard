import { CustomChartDonut, GreenChartDonut } from '@app/utils/ChartDonut';
import { help } from '@app/Github/PullRequests';
import React from 'react';
import { Card, CardTitle, CardBody } from '@patternfly/react-core';

export const CustomSLIChartDonutCard = (props) => {
    const title = props.title
    const donutChartColorScale = props.donutChartColorScale
    const data = props.data
    const type = props.type
    let red_count = 0
    let yellow_count = 0
    let green_count = 0
    let donutChartTitle = 0
    let donutChartData = [{}]
    let donutChartLegend = [{}]

    switch (donutChartColorScale.length) {
        case 2:
            red_count = data?.filter(x => x[type].signal == "red").length
            yellow_count = data?.filter(x => x[type].signal == "yellow").length

            donutChartTitle = yellow_count + red_count
            donutChartData = [{ x: 'Red', y: red_count }, { x: 'Yellow', y: yellow_count }]
            donutChartLegend = [{ name: 'Red: ' + red_count }, { name: 'Yellow: ' + yellow_count }]
            break;
        case 1:
            // response sli and component assignment sli
            red_count = data?.filter(x => x[type].signal == "red").length

            donutChartTitle = red_count
            donutChartData = [{ x: 'Red', y: red_count }]
            donutChartLegend = [{ name: 'Red: ' + red_count }]
            break;
        default:
            red_count = props.sli?.red_sli
            yellow_count = props.sli?.yellow_sli
            green_count = props.sli?.green_sli

            donutChartTitle = red_count + yellow_count + green_count
            donutChartData = [{ x: 'Red', y: red_count }, { x: 'Yellow', y: yellow_count }, { x: 'Green', y: green_count }]
            donutChartLegend = [{ name: 'Red: ' + red_count }, { name: 'Yellow: ' + yellow_count }, { name: 'Green: ' + green_count }]
            break;
    }

    const getHeight = () => {
        if (props.title == 'Global Bug SLI Status') {
            return 60
        }

        return 30
    }

    return (
        <Card style={{ width: '100%', height: '100%', textAlign: 'center' }}>
            <CardTitle>
                <div style={{ height: getHeight() }}>
                    {title}
                    {props.title == 'Resolution Time Bug SLI' && (
                        help(
                            <div><b>Number of Blocker, Critical, and Major bugs that meet Resolution Time Bug SLO.</b>
                                <br /><br /><b>Red Blocker Bug Resolution Time SLI:</b> Unresolved for more than 10 days.
                                <br /><br /><b>Yellow Blocker Bug Resolution Time SLI:</b> Unresolved for more than 5 days.
                                <br /><br /><b>Red Critical Bug Resolution Time SLI:</b> Unresolved for more than 20 days.
                                <br /><br /><b>Yellow Critical Bug Resolution Time SLI:</b> Unresolved for more than 10 days.
                                <br /><br /><b>Red Major Bug Resolution Time SLI:</b> Unresolved for more than 40 days.
                                <br /><br /><b>Yellow Major Bug Resolution Time SLI:</b> Unresolved for more than 20 days.
                            </div>)
                    )}
                    {props.title == 'Response Time Bug SLI' && (
                        help(<div><b>Number of Blocker and Critical bugs that meet Response Time Bug SLO.</b> <br /><br /><b>Red SLI</b>: Unassigned for more than 2 days on Blocker and Critical bugs.</div>)
                    )}
                    {props.title == 'Priority Triage Time Bug SLI' && (
                        help(<div><b>Number of untriaged bugs that meet Priority Triage Time Bug SLO.</b> <br /><br /><b>Red SLI</b>: Priority undefined for more than 2 days on untriaged bugs. <br /><br /> <b>Yellow SLI</b>: Priority undefined for more than 1 day on untriaged bugs.</div>)
                    )}
                    {props.title == 'Component Assignment Triage Time SLI' && (
                        help(<div><b>Number of bugs that meet Component Assignment Triage Time Bug SLO.</b> <br /><br /><b>Red SLI</b>: Component undefined for more than 1 day.</div>)
                    )}
                    {props.title == 'Global Bug SLI Status' && (
                        help(<div>Global Status of each opened bug.</div>)
                    )}
                </div>
            </CardTitle>
            {donutChartTitle > 0 && <CardBody>
                <CustomChartDonut
                    donutChartColorScale={donutChartColorScale}
                    donutChartData={donutChartData}
                    donutChartTitle={donutChartTitle}
                    donutChartLegend={donutChartLegend}
                >
                </CustomChartDonut>
            </CardBody>}
            {donutChartTitle == 0 && <CardBody>
                <GreenChartDonut></GreenChartDonut>
            </CardBody>}
        </Card>
    )
}