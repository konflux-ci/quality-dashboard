import React, { useLayoutEffect, useRef, useState } from 'react';
import { Card, CardBody, CardTitle, Popover, Title } from '@patternfly/react-core';
import {
  Chart,
  ChartAxis,
  ChartGroup,
  ChartLine,
  ChartThemeColor,
  createContainer,
} from '@patternfly/react-charts';
import { DashboardLineChartData } from '@app/utils/sharedComponents';
import { HelpIcon } from '@patternfly/react-icons';
import { getLabels } from '@app/utils/utils';

export interface PrsStatistics {
  metrics: Metrics[];
  summary: Summary;
}

export interface Summary {
  open_prs: number;
  merged_prs: number;
  merge_avg: string;
  retest_avg: number;
  retest_before_merge_avg: number;
}

export interface Metrics {
  date: string;
  created_prs_count: number;
  merged_prs_count: number;
}

export const help = (desc: string) => {
  return (
    <Popover
      bodyContent={
        <div>
          {desc}
        </div>
      }
    >
      <button
        type="button"
        aria-label="More info for name field"
        onClick={(e) => e.preventDefault()}
        aria-describedby="modal-with-form-form-name"
        className="pf-c-form__group-label-help"
        title=""
      >
        <HelpIcon style={{ marginLeft: 5 }} noVerticalAlign />
      </button>
    </Popover>
  );
}

export const PullRequestCard = (props) => {
  return (
    <Card style={{ width: '100%', height: '100%', textAlign: 'center' }}>
      <CardTitle>
        <div>{props.title}</div>
        <div style={{ color: 'grey', fontSize: 12 }}>
          {props.subtitle}
          {props.title == 'Retest Avg' && (
            help("Retests: calculate an average how many /test and /retest comments were in total issued for pull requests opened in selected time range")
          )}
          {props.title == 'Retest Before Merge Avg' && (
            help("Retests to merge: calculate an average how many /test and /retest comments were issued after the last code push")
          )}
          {props.title == 'Time To Merge PR Avg Days' && (
            help("Average time to merge a PR: calculate an average of how many days were needed to merge a PR (difference between creation and merged date)")
          )}
        </div>
      </CardTitle>
      {props.repo == undefined && (
        <CardBody>
          <div style={{ fontSize: 25 }}>{props.total}</div>
        </CardBody>
      )}
      {props.repo != undefined && (
        <CardBody>
          <div style={{ fontSize: 25 }}>
            {props.repo.coverage.average_to_retest_before_merge == 0
              ? 'N/A'
              : props.repo.coverage.average_to_retest_before_merge}
          </div>
        </CardBody>
      )}
    </Card>
  );
};

const getMaxY = (data: DashboardLineChartData) => {
  let maxY = 0;
  Object.keys(data).map((key, _) => {
    data[key].data.forEach((metric) => {
      if (metric.y > maxY) {
        maxY = metric.y;
      }
    });
  });
  return maxY;
};

export const PullRequestsGraphic = (props) => {
  const ZoomVoronoiContainer = createContainer("zoom", "voronoi");
  const legendData = [
    { childName: 'created prs', name: 'Created PRs' },
    { childName: 'merged prs', name: 'Merged PRs' },
    { childName: 'retest avg', name: 'Retest Avg' },
    { childName: 'retest before merge avg', name: 'Retest Before Merge Avg' },
  ];
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
    CREATED_PRS: { data: [] },
    MERGED_PRS: { data: [] },
    RETEST_AVG: { data: [] },
    RETEST_BEFORE_MERGE_AVG: { data: [] },
  };

  props.metrics.forEach((metric) => {
    beautifiedData['CREATED_PRS'].data.push({
      name: 'created_prs',
      x: new Date(metric.date).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: '2-digit' }),
      y: metric.created_prs_count,
    });
    beautifiedData['MERGED_PRS'].data.push({
      name: 'merged_prs',
      x: new Date(metric.date).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: '2-digit' }),
      y: metric.merged_prs_count,
    });
    beautifiedData['RETEST_AVG'].data.push({
      name: 'retest_avg',
      x: new Date(metric.date).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: '2-digit' }),
      y: props.summary.retest_avg > 0.01 ? metric.retest_avg : 0,
    });
    beautifiedData['RETEST_BEFORE_MERGE_AVG'].data.push({
      name: 'retest_before_merge_avg',
      x: new Date(metric.date).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: '2-digit' }),
      y: props.summary.retest_before_merge_avg > 0.01 ? metric.retest_before_merge_avg : 0,
    });
  });

  const maxY = getMaxY(beautifiedData);

  return (
    <div style={{ height: '100%', width: '100%', minHeight: '600px' }} className={'pf-c-card'} ref={ref}>
      <div style={{ height: height + 'px', width: width + 'px', background: 'white', textAlign: 'center' }}>
        <Title style={{ textAlign: 'center', marginLeft: 20, marginTop: 20 }} headingLevel={'h2'}>
          Pull Requests over the selected time range
        </Title>
        <Chart
          ariaDesc="Average number of pets"
          containerComponent={
            <ZoomVoronoiContainer
              labels={({ datum }) => getLabels(datum, "created_prs")}
              voronoiDimension="x"
              voronoiPadding={0}
              constrainToVisibleArea
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
            top: 80,
          }}
          themeColor={ChartThemeColor.green}
          width={width}
          height={height}
        >
          <ChartAxis fixLabelOverlap={true} showGrid />
          <ChartAxis fixLabelOverlap={true} dependentAxis showGrid />
          <ChartGroup>
            {Object.keys(beautifiedData).map((key, index) => {
              return (
                <ChartLine key={index} data={beautifiedData[key].data} name={key.toLowerCase().replace(/_/gi, ' ')} />
              );
            })}
          </ChartGroup>
        </Chart>
      </div>
    </div>
  );
};
