import React from 'react';
import { Chart, ChartAxis, ChartBar, ChartGroup, ChartStack, ChartThemeColor, ChartVoronoiContainer } from '@patternfly/react-charts';


const getData = (bugs, name, search) =>  {
    return bugs?.map((x) => {
        return {
            name: name,
            x: x.component,
            y: bugs?.filter(y => y.component == x.component && y[search]?.signal != "green").length,
            red: bugs?.filter(y => y.component == x.component && y[search]?.signal == "red").length,
            yellow: bugs?.filter(y => y.component == x.component && y[search]?.signal == "yellow").length,
        };
    }).filter(
        (elem, index, arr) => index === arr.findIndex((t) => t.x === elem.x)
    );
}


export const SLIsStackChart = (props) => {
    const bugs = props.bugSLIs.bugs
    const resolution = getData(bugs, "Bugs Meeting Resolution Time Bug SLO", "resolution_sli")
    const response = getData(bugs, "Bugs Meeting Response Time Bug SLO", "response_sli")
    const triage = getData(bugs, "Bugs Meeting Triage Time Bug SLO", "triage_sli")

    return <div>
        {(resolution.length > 0 || response.length > 0 || triage.length > 0) && <div style={{ margin: 'auto', height: '275px', width: '800px' }}>
            <Chart
                containerComponent={<ChartVoronoiContainer labels={({ datum }) => `${datum.name}: ${datum.y} (Red: ${datum.red}, Yellow: ${datum.yellow})`} constrainToVisibleArea />}
                legendData={[{ name: 'Bugs Meeting Resolution Time Bug SLO' }, { name: 'Bugs Meeting Response Time Bug SLO' }, { name: 'Bugs Meeting Triage Time Bug SLO' }]}
                legendOrientation="horizontal"
                legendPosition="bottom"
                height={275}
                name="chart3"
                padding={{
                    bottom: 75, // Adjusted to accommodate legend
                    left: 50,
                    right: 50,
                    top: 50
                }}
                // colorScale={["#ea2745", "#fbe424"]}
                width={800}
            >
                <ChartAxis />
                <ChartAxis dependentAxis showGrid />
                <ChartStack>
                    <ChartBar data={resolution} />
                    <ChartBar data={response} />
                    <ChartBar data={triage} />
                </ChartStack>
            </Chart>
        </div>
        }
    </div>

}