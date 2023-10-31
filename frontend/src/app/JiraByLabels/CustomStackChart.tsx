import React from 'react';
import { Chart, ChartAxis, ChartBar, ChartStack, ChartThemeColor, ChartTooltip } from '@patternfly/react-charts';

export const CustomStackChart = (props) => {
    const data = props.data
    const legend = props.legend

    return (
        <div style={{ margin: 'auto', height: '300px', width: '400px' }}>
            <Chart
                legendData={legend}
                legendPosition="bottom"
                legendAllowWrap={true}
                height={270}
                name="chart3"
                padding={{
                    bottom: 75, // Adjusted to accommodate legend
                    left: 120,
                    right: 50,
                    top: 50
                }}
                themeColor={ChartThemeColor.multiOrdered}
                width={400}
            >
                <ChartAxis />
                <ChartAxis dependentAxis showGrid />
                <ChartStack horizontal>
                    {data.map((d, index) => {
                        return (
                            <ChartBar
                                key={index}
                                data={d}
                                labelComponent={<ChartTooltip constrainToVisibleArea />}
                            />
                        )
                    })}
                </ChartStack>
            </Chart>
        </div>
    )
} 