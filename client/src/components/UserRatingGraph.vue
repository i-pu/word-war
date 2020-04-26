<script lang="ts">
import Vue, { PropType } from 'vue'
import { ChartData, ChartOptions } from 'chart.js'
import { Line } from 'vue-chartjs'

export default {
  name: 'UserRatingGraph',
  extends: Line,
  props: {
    ratings: {
      type: Array,
      default: () => []
    }
  },
  data (): {
    options: ChartOptions,
  } {
    return {
      options: {
        legend: {
          display: false,
        },
        scales: {
          xAxes: [{
            type: 'time',
            distribution: 'series',
            time: {
              unit: 'day'
            }
          }]
        },
        elements: {
          line: {
            tension: 0
          }
        }
      }
    }
  },
  computed: {
    data(): ChartData {
      return {
        datasets: [
          {
            // @ts-ignore
            data: this.ratings.map(({ rating, date }) => ({ t: date, y: rating })),
            fill: false,
          }
        ]
      }
    }
  },
  mounted () {
    // @ts-ignore
    console.log(this.ratings);
    // @ts-ignore
    (this as any).renderChart(this.data, this.options)
  },
}
</script>