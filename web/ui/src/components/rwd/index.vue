<template lang="pug">
el-row(:gutter="10")
  el-col.offset(:xs="1", :sm="1", :md="2", :lg="2")
    | &nbsp;
  el-col(:xs="22", :sm="22", :md="20", :lg="20")
    slot
  el-col.offset(:xs="1", :sm="1", :md="2", :lg="2")
    | &nbsp;
</template>