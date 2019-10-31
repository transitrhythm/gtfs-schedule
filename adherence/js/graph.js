    var lastDate = 0;
    var data = []

    function getDayWiseTimeSeries(baseval, count, yrange) {
      var i = 0;
      while (i < count) {
        var x = baseval;
        var y = Math.floor(Math.random() * (yrange.max - yrange.min + 1)) + yrange.min;

        data.push({
          x,
          y
        });
        lastDate = baseval
        baseval += 86400000;
        i++;
      }
    }

    getDayWiseTimeSeries(new Date('11 Feb 2017 GMT').getTime(), 10, {
      min: -10,
      max: 30
    })

    function getNewSeries(baseval, yrange) {
      var newDate = baseval + 86400000;
      lastDate = newDate
      data.push({
        x: newDate,
        y: Math.floor(Math.random() * (yrange.max - yrange.min + 1)) + yrange.min
      })
    }

    function resetData() {
      data = data.slice(data.length - 10, data.length);
    }

    new Vue({
      el: '#chart',
      components: {
        apexchart: VueApexCharts,
      },
      data: {
        series: [{
          data: data.slice()
        }],
        chartOptions: {
          chart: {
            animations: {
              enabled: true,
              easing: 'linear',
              dynamicAnimation: {
                speed: 1000
              }
            },
            toolbar: {
                autoSelected: 'zoom'
            },
            zoom: {
                type: 'x',
                enabled: true
            }
          },
          dataLabels: {
            enabled: false
          },
          stroke: {
            curve: 'smooth'
          },

          title: {
            text: 'Schedule Adherence Route: ### Trip: #####',
            align: 'left'
          },
          markers: {
            size: 1
          },
          xaxis: {
            type: 'datetime',
            range: 777600000,
          },
          yaxis: {
            max: 30
          },
          legend: {
            show: false
          }
        }
      },
      mounted: function () {
        this.intervals()
      },
      methods: {
        intervals: function () {
          var me = this
          window.setInterval(function () {
            getNewSeries(lastDate, {
              min: -10,
              max: 30
            })
            
            me.$refs.realtimeChart.updateSeries([{
              data: data
            }])
          }, 15000)

          // every 60 seconds, we reset the data to prevent memory leaks
          window.setInterval(function () {
            resetData()
            me.$refs.realtimeChart.updateSeries([{
              data
            }], false, true)
          }, 60000)
        }
      }
    })
