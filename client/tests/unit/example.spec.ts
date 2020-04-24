import { shallowMount } from '@vue/test-utils'
import Top from '@/views/Top.vue'

describe('Top.vue', () => {
  it('renders props.msg when passed', () => {
    const msg = 'Jest'
    const wrapper = shallowMount({
      template: '<p>Hello {{ msg }}</p>'
    }, {
      propsData: { msg },
    })
    expect(wrapper.text()).toMatch(`Hello ${msg}`)
  })
})
