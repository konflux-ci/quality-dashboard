import React from 'react';
import { Chart, ChartAxis, ChartGroup, ChartLine, ChartThemeColor, ChartLegendTooltip, createContainer } from '@patternfly/react-charts';

class FakeChart extends React.Component {
  render() {
    // Note: Container order is important
    const CursorVoronoiContainer = createContainer("voronoi", "cursor");
    const legendData = [{ childName: 'cats', name: 'Cats' }, { childName: 'dogs', name: 'Dogs', symbol: { type: 'dash' }}, { childName: 'birds', name: 'Birds' }, { childName: 'mice', name: 'Mice' }];

    return (
      <div style={{ height: '275px', width: '450px' }}>
        <Chart
          ariaDesc="Average number of pets"
          ariaTitle="Line chart example"
          containerComponent={
            <CursorVoronoiContainer
              cursorDimension="x"
              labels={({ datum }) => `${datum.y}`}
              labelComponent={<ChartLegendTooltip legendData={legendData} title={(datum) => datum.x}/>}
              mouseFollowTooltips
              voronoiDimension="x"
              voronoiPadding={50}
            />
          }
          legendData={legendData}
          legendPosition="bottom"
          height={275}
          maxDomain={{y: 10}}
          minDomain={{y: 0}}
          name="chart2"
          padding={{
            bottom: 75, // Adjusted to accommodate legend
            left: 50,
            right: 50,
            top: 50
          }}
          themeColor={ChartThemeColor.green}
          width={450}
        >
          <ChartAxis tickValues={[2, 3, 4]} />
          <ChartAxis dependentAxis showGrid tickValues={[2, 5, 8]} />
          <ChartGroup>
            <ChartLine
              data={[
                { x: '2015', y: 1 },
                { x: '2016', y: 2 },
                { x: '2017', y: 5 },
                { x: '2018', y: 3 }
              ]}
              name="cats"
            />
            <ChartLine
              data={[
                { x: '2015', y: 2 },
                { x: '2016', y: 1 },
                { x: '2017', y: 7 },
                { x: '2018', y: 4 }
              ]}
              name="dogs"
              style={{
                data: {
                  strokeDasharray: '3,3'
                }
              }}
            />
            <ChartLine
              data={[
                { x: '2015', y: 3 },
                { x: '2016', y: 4 },
                { x: '2017', y: 9 },
                { x: '2018', y: 5 }
              ]}
              name="birds"
            />
            <ChartLine
              data={[
                { x: '2015', y: 3 },
                { x: '2016', y: 3 },
                { x: '2017', y: 8 },
                { x: '2018', y: 7 }
              ]}
              name="mice"
            />
          </ChartGroup>
        </Chart>
      </div>
    );
  }
}