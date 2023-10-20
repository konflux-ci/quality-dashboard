import React from 'react';
import { Chart, ChartAxis, ChartBar, ChartGroup, ChartStack, ChartThemeColor, ChartVoronoiContainer } from '@patternfly/react-charts';


const getData = (bugs, name, sli) =>  {
    return bugs?.map((x) => {
        return {
            name: name,
            x: x.component,
            y: bugs?.filter(y => y.component == x.component && y.global_sli == sli).length
        };
    }).filter(
        (elem, index, arr) => index === arr.findIndex((t) => t.x === elem.x)
    );
}


export const SLIsStackChart = (props) => {
    const bugSLIs = props.bugSLIs.bugs
    const red = getData(bugSLIs, "Red", "red")
    const yellow = getData(bugSLIs, "Yellow", "yellow")

    return <div>
        {(red.length > 0 || yellow.length > 0) && <div style={{ margin: 'auto', height: '275px', width: '800px' }}>
            <Chart
                containerComponent={<ChartVoronoiContainer labels={({ datum }) => `${datum.name}: ${datum.y}`} constrainToVisibleArea />}
                legendData={[{ name: 'Red' }, { name: 'Yellow' }]}
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
                colorScale={["#ea2745", "#fbe424"]}
                width={800}
            >
                <ChartAxis />
                <ChartAxis dependentAxis showGrid />
                <ChartStack>
                    <ChartBar data={red} />
                    <ChartBar data={yellow} />
                </ChartStack>
            </Chart>
        </div>
        }
    </div>

}