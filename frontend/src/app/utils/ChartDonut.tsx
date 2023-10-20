import React from 'react';
import { ChartDonut } from '@patternfly/react-charts';

export const CustomChartDonut = (props) => {
    const colorScale = props.donutChartColorScale
    const data = props.donutChartData
    const legendData = props.donutChartLegend
    const title = props.donutChartTitle

    return <div>
        <ChartDonut
            colorScale={colorScale}
            constrainToVisibleArea
            data={data}
            labels={({ datum }) => `${datum.x}: ${datum.y}`}
            legendData={legendData}
            legendPosition="bottom"
            name="chart10"
            padding={{
                bottom: 70, // Adjusted to accommodate legend
                left: 20,
                right: 50, // Adjusted to accommodate subTitle
                top: 20
            }}
            width={350}
            title={title}
        />
    </div>
}

export const GreenChartDonut = (props) => {

    return <div>
        <ChartDonut
            colorScale={["#ea2745", "#fbe424", "#61ad50"]}
            constrainToVisibleArea
            data={[{ x: '', y: 0 }, { x: '', y: 0 }, { x: 'All green', y: 1 }]}
            legendData={[{ name: 'Red: 0' }, { name: 'Yellow: 0' }]}
            legendPosition="bottom"
            name="chart10"
            padding={{
                bottom: 70, // Adjusted to accommodate legend
                left: 20,
                right: 50, // Adjusted to accommodate subTitle
                top: 20
            }}
            width={350}
            title={"0"}
        />
    </div>
}