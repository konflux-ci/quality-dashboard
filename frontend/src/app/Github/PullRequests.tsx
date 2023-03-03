import React, { useLayoutEffect, useRef, useState } from "react";
import { Card, CardBody, CardTitle, Title } from "@patternfly/react-core";
import { Chart, ChartAxis, ChartGroup, ChartLegendTooltip, ChartLine, ChartThemeColor, createContainer } from "@patternfly/react-charts";
import { DashboardLineChartData } from "@app/utils/sharedComponents";

export interface PrsStatistics {
    metrics: (Metrics)[];
    summary: Summary;
}

export interface Summary {
    open_prs: number;
    merged_prs: number;
    merge_avg: string;
}

export interface Metrics {
    date: string;
    created_prs_count: number;
    merged_prs_count: number;
}

export const PullRequestCard = (props) => {

    return (
        <Card style={{ width: "100%", height: "100%" }}>
            <CardTitle>
                {props.title}
            </CardTitle>
            <CardBody style={{ textAlign: 'center' }}>
                {props.total}
            </CardBody>
        </Card>
    );
}

const getMaxY = (data: DashboardLineChartData) => {
    let maxY = 0
    Object.keys(data).map((key, _) => {
        data[key].data.forEach((metric) => {
            if (metric.y > maxY) {
                maxY = metric.y
            }
        })
    })
    return maxY
}

export const PullRequestsGraphic = (props) => {
    const CursorVoronoiContainer = createContainer("voronoi", "cursor");
    const legendData = [{ childName: 'created prs', name: 'Created PRs' }, { childName: 'merged prs', name: 'Merged PRs' }];
    const ref = useRef<HTMLDivElement>(null);
    const [width, setWidth] = useState(0);
    const [height, setHeight] = useState(0);

    useLayoutEffect(() => {
        if (ref.current) {
            setWidth(ref.current.offsetWidth - 0);
            setHeight(ref.current.offsetHeight - 40);
        }
    }, []);

    // Prepare data for the line chart
    let beautifiedData: DashboardLineChartData = {
        "CREATED_PRS": { data: [] },
        "MERGED_PRS": { data: [] },
    };

    props.metrics.forEach(metric => {
        beautifiedData["CREATED_PRS"].data.push({ name: 'created_prs', x: new Date(metric.date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: metric.created_prs_count })
        beautifiedData["MERGED_PRS"].data.push({ name: 'merged_prs', x: new Date(metric.date).toLocaleDateString("en-US", { day: 'numeric', month: 'short' }), y: metric.merged_prs_count })
    });

    const maxY = getMaxY(beautifiedData)

    return (
        <div style={{ height: '100%', width: '100%', minHeight: "600px" }} className={"pf-c-card"} ref={ref}>
            <div style={{ height: height + 'px', width: width + 'px', background: "white", textAlign: "center" }}>
                <Title style={{ textAlign: "left", marginLeft: 20, marginTop: 20 }} headingLevel={'h2'}>Pull Requests over time</Title>
                <Chart
                    ariaDesc="Average number of pets"
                    ariaTitle="Line chart example"
                    containerComponent={
                        <CursorVoronoiContainer
                            cursorDimension="x"
                            labels={({ datum }) => `${datum.y}`}
                            labelComponent={<ChartLegendTooltip legendData={legendData} title={(datum) => datum.x} />}
                            mouseFollowTooltips
                            voronoiDimension="x"
                            voronoiPadding={50}
                        />
                    }
                    legendData={legendData}
                    legendPosition="bottom"
                    maxDomain={{ y: maxY + 1 }}
                    minDomain={{ y: 0 }}
                    name="chart2"
                    padding={{
                        bottom: 110,
                        left: 80,
                        right: 80,
                        top: 80
                    }}
                    themeColor={ChartThemeColor.green}
                    width={width}
                    height={height}
                >
                    <ChartAxis fixLabelOverlap={true} showGrid />
                    <ChartAxis fixLabelOverlap={true} dependentAxis showGrid />
                    <ChartGroup>
                        {
                            Object.keys(beautifiedData).map((key, index) => {
                                return <ChartLine
                                    key={index}
                                    data={beautifiedData[key].data}
                                    name={key.toLowerCase().replace(/_/gi, " ")}
                                />
                            })
                        }
                    </ChartGroup>
                </Chart>
            </div>
        </div>
    );
}