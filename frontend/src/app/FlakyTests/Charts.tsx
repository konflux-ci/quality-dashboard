import React, { useState, useRef, useLayoutEffect } from 'react';
import { Chart, ChartAxis, ChartLine, ChartArea, ChartScatter, ChartVoronoiContainer } from '@patternfly/react-charts';
import { Skeleton } from '@patternfly/react-core';

export const ImpactChart: React.FunctionComponent<{ data, x, y, secondaryData?}> = ({ data, x, y, secondaryData }) => {
    const ref = useRef<HTMLDivElement>(null);
    const [width, setWidth] = useState(100);
    const [height, setHeight] = useState(100);

    useLayoutEffect(() => {
        if (ref.current && ref.current.offsetWidth > 0) {
            setWidth(ref.current.offsetWidth * 0.8 - 10);
            setHeight(ref.current.offsetWidth * 0.8 * 0.4 - 20);
        }
    }, []);

    return (
        <div style={{ width: '100%', height: '100%', boxShadow: "none" }} className={"pf-c-card"} ref={ref}>
            <div style={{ height: height + 'px', width: width + 'px', background: "white", boxShadow: "none" }}>
                {
                    data && data.length > 0 && <Chart
                        ariaDesc="Global impact"
                        ariaTitle="Global Impact"
                        containerComponent={<ChartVoronoiContainer labels={({ datum }) => `${datum.name}: ${datum.y}`} constrainToVisibleArea />}
                        domain={{ y: [-1, 101] }}
                        domainPadding={{ x: 0, y: 0 }}
                        legendOrientation="vertical"
                        legendPosition="right"
                        legendData={[{ name: "Global Impact", symbol: { fill: "green" } }, { name: "Flaky test impact", symbol: { fill: "#6495ED" } }]}
                        height={height}
                        width={width}
                        name="chart1"
                        padding={{
                            bottom: 100,
                            left: 60,
                            right: 250,
                            top: 50
                        }}
                    >
                        <ChartAxis showGrid style={{ tickLabels: { angle: 0, fontSize: 9 } }} />
                        <ChartAxis showGrid dependentAxis />

                        <ChartLine data={data.map((datum) => { return { "name": datum[x], "x": datum[x], "y": datum[y] ? parseFloat(datum[y]) : 0 } })} />

                        {
                            secondaryData && <ChartLine style={{
                                data: { stroke: "green", strokeWidth: 2 }
                            }} data={secondaryData.map((datum) => { return { "name": datum[x], "x": datum[x], "y": datum[y] ? parseFloat(datum[y]) : 0 } })} />
                        }

                        <ChartArea style={{
                            data: {
                                fill: "#6495ED", fillOpacity: 0.3
                            }
                        }} data={data.map((datum) => { return { "name": datum[x], "x": datum[x], "y": datum[y] ? parseFloat(datum[y]) : 0 } })} />

                    </Chart>
                }
                {(!data || data.length == 0) &&
                    <div style={{ height: '100%', display: 'flex', alignItems: 'flex-end', justifyContent: 'space-between' }}>
                        <Skeleton height="25%" width="15%" screenreaderText="Loading contents" />
                        <Skeleton height="33%" width="15%" />
                        <Skeleton height="50%" width="15%" />
                        <Skeleton height="66%" width="15%" />
                        <Skeleton height="75%" width="15%" />
                        <Skeleton height="100%" width="15%" />
                    </div>
                }
            </div>
        </div>
    )
}

export const RegressionChart: React.FunctionComponent<{ data, x, y }> = ({ data, x, y }) => {
    const ref = useRef<HTMLDivElement>(null);
    const [width, setWidth] = useState(100);
    const [height, setHeight] = useState(100);

    useLayoutEffect(() => {
        if (ref.current && ref.current.offsetWidth > 0) {
            setWidth(ref.current.offsetWidth * 0.8 - 10);
            setHeight(ref.current.offsetWidth * 0.8 * 0.4 - 20);
        }
    }, []);

    return (
        <div style={{ width: '100%', height: '100%', boxShadow: "none" }} className={"pf-c-card"} ref={ref}>
            <div style={{ height: height + 'px', width: width + 'px', background: "white", boxShadow: "none" }}>
                {
                    data && data.length > 0 && <Chart
                        ariaDesc="Regression"
                        ariaTitle="Regression"
                        containerComponent={<ChartVoronoiContainer labels={({ datum }) => `${datum.name}: ${datum.y}`} constrainToVisibleArea />}
                        domain={{ y: [0, Math.max(...data.map(o => o.global_impact))], x: [0, Math.max(...data.map(o => o.jobs_executed))] }}
                        legendData={[{ name: "X-axis: Job count", symbol: { fill: "white" } }, { name: "Y-axis: % failed jobs", symbol: { fill: "white" } }, { name: "Job-Impact", symbol: { fill: "orange" } }, { name: "Regression", symbol: { fill: "darkgray" } }]}
                        domainPadding={{ x: 0, y: 0 }}
                        legendOrientation="vertical"
                        legendPosition="right"
                        height={height}
                        width={width}
                        name="chart1"
                        padding={{
                            bottom: 100,
                            left: 60,
                            right: 250,
                            top: 50
                        }}
                    >
                        <ChartAxis showGrid style={{ tickLabels: { angle: 0, fontSize: 9 } }} />
                        <ChartAxis showGrid dependentAxis />

                        <ChartScatter style={{ data: { fill: "orange" } }} data={data.filter(d => d.jobs_executed != 0).sort((aValue, bValue) => { return (aValue.jobs_executed as number) - (bValue.jobs_executed as number) }).map((datum) => { return { "name": datum[x], "x": datum[x], "y": datum[y] ? parseFloat(datum[y]) : 0 } })} />

                        <ChartLine style={{ data: { stroke: "darkgray" } }} data={data.filter(d => d.jobs_executed != 0).sort((aValue, bValue) => { return (aValue.jobs_executed as number) - (bValue.jobs_executed as number) }).map((datum) => { return { "name": "Regression", "x": datum[x], "y": datum.regression ? parseFloat(datum.regression) : 0 } })} />

                    </Chart>
                }
                {(!data || data.length == 0) &&
                    <div style={{ height: '100%', display: 'flex', alignItems: 'flex-end', justifyContent: 'space-between' }}>
                        <Skeleton height="25%" width="15%" screenreaderText="Loading contents" />
                        <Skeleton height="33%" width="15%" />
                        <Skeleton height="50%" width="15%" />
                        <Skeleton height="66%" width="15%" />
                        <Skeleton height="75%" width="15%" />
                        <Skeleton height="100%" width="15%" />
                    </div>
                }
            </div>
        </div>
    )
}