<template>
  <div class="hello">
  </div>
</template>

<script>
import { EventBus } from '../event-bus.js'
import Events from '../store/event-api.js'

export default {
  name: 'HelloWorld',
  created () {
    EventBus.$emit(Events.toggleDrawer)
  }
}
</script>


