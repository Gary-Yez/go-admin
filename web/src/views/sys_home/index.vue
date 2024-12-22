<template>
  <div>
    <el-row :gutter="15">
      <el-col :span="6">
        <el-card class="container" shadow="never">
          <el-statistic title="协程数量" :value="pageData.go_threads">
            <template #title>
              <div class="flex items-center">
                <el-icon class="mr-[5px]" size="20" color="#0099ff">
                  <component is="Cpu"></component>
                </el-icon>
                <strong class="text-[14px]">协程数量</strong>
              </div>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="container" shadow="never">
          <el-statistic title="占用内存" :value="(pageData.use_memory / 1024 / 1024).toFixed(2)">
            <template #title>
              <div class="flex items-center">
                <el-icon class="mr-[5px]" size="20" color="#0099ff">
                  <component is="Coin"></component>
                </el-icon>
                <strong class="text-[14px]">占用内存</strong>
              </div>
            </template>
            <template #suffix>M</template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="container" shadow="never">
          <el-statistic title="运行时长" :value="formatTimeDifference(+new Date(),pageData.start_time)">
            <template #title>
              <div class="flex items-center">
                <el-icon class="mr-[5px]" size="20" color="#0099ff">
                  <component is="Timer"></component>
                </el-icon>
                <strong class="text-[14px]">运行时长</strong>
              </div>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="container" shadow="never">
          <el-statistic title="启动时间" :value="formatTime(pageData.start_time)">
            <template #title>
              <div class="flex items-center">
                <el-icon class="mr-[5px]" size="20" color="#0099ff">
                  <component is="Timer"></component>
                </el-icon>
                <strong class="text-[14px]">启动时间</strong>
              </div>
            </template>
          </el-statistic>
        </el-card>
      </el-col>
    </el-row>
    <el-card class="container mt-[15px]" header="系统信息" shadow="never">
      <div class="sys_info">
        <div class="system_info_item">
          <div class="title">系统负载</div>
          <el-progress type="dashboard" :stroke-width="8" :percentage="pageData.loadAvg?.load1.toFixed(2)" striped  striped-flow>
            <template #default="{ percentage }">
              <span class="percentage-value font-bold">{{ percentage }}%</span>
            </template>
          </el-progress>
          <div class="detail_info">运行流畅</div>
        </div>
        <div class="system_info_item">
          <div class="title">CPU占用率</div>
          <el-progress type="dashboard" :stroke-width="8" :percentage="cpuPercent" striped  striped-flow>
            <template #default="{ percentage }">
              <span class="percentage-value font-bold">{{ percentage }}%</span>
            </template>
          </el-progress>
          <div class="detail_info">{{ pageData.cpu_status?.counts }} 核心</div>
        </div>
        <div class="system_info_item">
          <div class="title">内存占用率</div>
          <el-progress type="dashboard" :stroke-width="8" :percentage="pageData.mem_status?.percent" striped  striped-flow>
            <template #default="{ percentage }">
              <span class="percentage-value font-bold">{{ percentage }}%</span>
            </template>
          </el-progress>
          <div class="detail_info">{{ pageData.mem_status?.used.toFixed(0) }} / {{ pageData.mem_status?.total.toFixed(0) }}(MB)</div>
        </div>
        <div v-for="item in pageData.disks" class="system_info_item">
          <div class="title">磁盘{{ item.name }}</div>
          <el-progress type="dashboard" :stroke-width="8" :percentage="item?.percent.toFixed(2)" striped  striped-flow>
            <template #default="{ percentage }">
              <span class="percentage-value font-bold">{{ percentage }}%</span>
            </template>
          </el-progress>
          <div class="detail_info">{{ (item.used / 1024).toFixed(0) }}G / {{ (item.total / 1024).toFixed(0) }}G</div>
        </div>
      </div>
    </el-card>
    <el-row :gutter="15">
      <el-col :span="12">
        <el-card class="container mt-[15px]" shadow="never">
          <template #header>
            <div class="flex items-center justify-between">
              <div>网络IO / KB每秒</div>
              <div class="w-[150px]">
                <el-select v-model="netName">
                  <el-option label="全部" value="all"></el-option>
                  <el-option v-for="item in pageData.netIO" :label="item.name" :value="item.name"></el-option>
                </el-select>
              </div>
            </div>
          </template>
          <div class="ioBox">
            <div>
              <div class="label">
                <div class="dot bg-amber-300"></div>
                <span>上行</span>
              </div>
              <div class="value">{{ currentNetIo.bytesSent.toFixed(2) }} KB</div>
            </div>
            <div>
              <div class="label">
                <div class="dot bg-blue-300"></div>
                <span>下行</span>
              </div>
              <div class="value">{{ currentNetIo.bytesRecv.toFixed(2) }} KB</div>
            </div>
            <div>
              <div class="label">总发送</div>
              <div class="value">{{ currentNetIo.allSend.toFixed(2) }} GB</div>
            </div>
            <div>
              <div class="label">总接收</div>
              <div class="value">{{ currentNetIo.allRecv.toFixed(2) }} GB</div>
            </div>
          </div>
          <div id="netIo" class="h-[400px]"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="container mt-[15px]" shadow="never">
          <template #header>
            <div class="flex items-center justify-between">
              <div>磁盘IO / MB每秒</div>
              <div class="w-[150px]">
                <el-select v-model="diskName">
                  <el-option label="全部" value="all"></el-option>
                  <el-option v-for="item in pageData.diskIO" :label="item.name" :value="item.name"></el-option>
                </el-select>
              </div>
            </div>
          </template>
          <div class="ioBox">
            <div>
              <div class="label">
                <div class="dot bg-amber-300"></div>
                <span>读取</span>
              </div>
              <div class="value">{{ currentdiskIo.readBytes.toFixed(2) }} MB</div>
            </div>
            <div>
              <div class="label">
                <div class="dot bg-blue-300"></div>
                <span>写入</span>
              </div>
              <div class="value">{{ currentdiskIo.writeBytes.toFixed(2) }} MB</div>
            </div>
            <div>
              <div class="label">每秒读写</div>
              <div class="value">{{ currentdiskIo.count }}</div>
            </div>
            <div>
              <div class="label">IO延迟</div>
              <div class="value">{{ currentdiskIo.ioTime }} ms</div>
            </div>
          </div>
          <div id="diskIo" class="h-[400px]"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
  import * as echarts from 'echarts';
  import {SysHomeApi} from "../../apis/sys_home.ts";
  import {computed, onBeforeUnmount, onMounted, ref} from "vue";
  import {formatTime} from "../../utils/formatTime.ts";
  const pageData:any = ref({})
  const oldData:any = ref([])
  let jsq:any
  let netIoChart:any
  let diskIoChart:any
  const netName = ref("all")
  const diskName = ref("all")

  const ioDataList = computed(() => {
    return oldData.value.map((item:any) => {
      if (netName.value === 'all'){
        return item.netIO.reduce((acc:any, cur:any) => {
          acc.bytesSent += cur.bytesSent;
          acc.bytesRecv += cur.bytesRecv;
          return acc
        },{
          bytesSent:0,
          bytesRecv:0,
        })
      }else{
        return item.netIO.find((io:any)=>io.name == netName.value)
      }
    })
  })
  const currentNetIo = computed(()=>{
    if (ioDataList.value.length >= 2){
      return {
        bytesSent:(ioDataList.value[ioDataList.value.length-1].bytesSent - ioDataList.value[ioDataList.value.length-2].bytesSent) / 1024,
        bytesRecv:(ioDataList.value[ioDataList.value.length-1].bytesRecv - ioDataList.value[ioDataList.value.length-2].bytesRecv) /1024,
        allSend:(ioDataList.value[ioDataList.value.length-1]?.bytesSent)/1024/1024/1024,
        allRecv:(ioDataList.value[ioDataList.value.length-1]?.bytesRecv)/1024/1024/1024,
      }
    }else{
      return {
        bytesSent:0,
        bytesRecv:0,
        allSend:(ioDataList.value[ioDataList.value.length-1]?.bytesSent)/1024/1024/1024 || 0,
        allRecv:(ioDataList.value[ioDataList.value.length-1]?.bytesRecv)/1024/1024/1024 || 0,
      }
    }
  })
  const diskDataList = computed(() => {
    return oldData.value.map((item:any) => {
      if (diskName.value === 'all'){
        let data:any = Object.values(item.diskIO).reduce((acc:any, cur:any) => {
          acc.readBytes += cur.readBytes;
          acc.readCount += cur.readCount;
          acc.writeBytes += cur.writeBytes;
          acc.writeCount += cur.writeCount;
          acc.ioTime += cur.ioTime ;
          return acc
        },{
          readBytes:0,
          readCount:0,
          writeBytes:0,
          writeCount:0,
          ioTime:0,
        })
        data.ioTime = data.ioTime / Object.values(item.diskIO).length;
        return data
      }else{
        return item.diskIO[diskName.value]
      }
    })
  })
  const currentdiskIo = computed(()=>{
    if (diskDataList.value.length >= 2){
      let current = diskDataList.value[diskDataList.value.length-1]
      let last = diskDataList.value[diskDataList.value.length-2]
      let count = (current.readCount + current.writeCount) - (last.readCount + last.writeCount)
      console.log(current)
      return {
        readBytes:(current.readBytes - last.readBytes) / 1024 /1024,
        writeBytes:(current.writeBytes - last.writeBytes) /1024 /1024,
        count:count,
        // allRecv:(diskDataList.value[diskDataList.value.length-1]?.bytesRecv)/1024/1024/1024,
        ioTime:current.ioTime
      }
    }else{
      return {
        readBytes:0,
        writeBytes:0,
        count:0,
        ioTime:0
        // allSend:(diskDataList.value[diskDataList.value.length-1]?.bytesSent)/1024/1024/1024 || 0,
        // allRecv:(diskDataList.value[diskDataList.value.length-1]?.bytesRecv)/1024/1024/1024 || 0,
      }
    }
  })



  const buildIoOption = ()=>{
    let times = []
    let data = []
    let data2 = []
    for (let i = 1; i < oldData.value.length; i++) {
      times.push(oldData.value[i].time)
      if (ioDataList.value[i - 1]){
        data.push((ioDataList.value[i].bytesSent - ioDataList.value[i - 1].bytesSent) / 1024)
        data2.push((ioDataList.value[i].bytesRecv - ioDataList.value[i - 1].bytesRecv) / 1024)
      }else{
        data.push(0)
        data2.push(0)
      }
    }
    return {
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: times
      },
      yAxis: {
        type: 'value'
      },
      grid: {
        left: '3%',
        right: '3%',
        bottom: '3%',
        containLabel: true
      },
      series: [
        {
          data: data,
          type: 'line',
          stack: 'Total',
          smooth:true,
          lineStyle: {
            width: 0
          },
          showSymbol: false,
          areaStyle: {
            opacity: 0.5,
            color: "#f7b851"
          }
        },
        {
          data: data2,
          type: 'line',
          stack: 'Total',
          smooth:true,
          lineStyle: {
            width: 0
          },
          showSymbol: false,
          areaStyle: {
            opacity: 0.5,
            color: "#52a9ff"
          }
        }
      ]
    }
  }
  const buildDiskIoOption = ()=>{
    let times = []
    let data = []
    let data2 = []
    for (let i = 1; i < oldData.value.length; i++) {
      times.push(oldData.value[i].time)
      if (diskDataList.value[i - 1]){
        data.push((diskDataList.value[i].writeBytes - diskDataList.value[i - 1].writeBytes) / 1024 / 1024)
        data2.push((diskDataList.value[i].readBytes - diskDataList.value[i - 1].readBytes) / 1024 / 1024)
      }else{
        data.push(0)
        data2.push(0)
      }
    }
    return {
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: times
      },
      yAxis: {
        type: 'value'
      },
      grid: {
        left: '3%',
        right: '3%',
        bottom: '3%',
        containLabel: true
      },
      series: [
        {
          data: data,
          type: 'line',
          stack: 'Total',
          smooth:true,
          lineStyle: {
            width: 0
          },
          showSymbol: false,
          areaStyle: {
            opacity: 0.5,
            color: "#6cc0cf"
          }
        },{
          data: data2,
          type: 'line',
          stack: 'Total',
          smooth:true,
          lineStyle: {
            width: 0
          },
          showSymbol: false,

          areaStyle: {
            opacity: 0.5,
            color: "#ff4683"
          }
        }
      ]
    }
  }

  const getPageData = async () => {
    try {
      const response = await SysHomeApi.Statistic()
      pageData.value = response.data
      if (oldData.value.length>= 7) {
        oldData.value.shift();
      }
      let now = new Date().toLocaleTimeString();
      oldData.value.push({
        time:now,
        ...response.data
      });
      netIoChart.setOption(buildIoOption());
      diskIoChart.setOption(buildDiskIoOption());
    }catch (e) {
      console.log(e)
    }
    jsq = setTimeout(getPageData,2000)
  }
  const formatTimeDifference = (timestamp1:number, timestamp2:number) => {
    // 计算两个时间戳的差值（以毫秒为单位）
    let diff = Math.abs(timestamp2 - timestamp1);

    // 将差值转换为天数、小时数和分钟数
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    diff -= days * (1000 * 60 * 60 * 24);

    const hours = Math.floor(diff / (1000 * 60 * 60));
    diff -= hours * (1000 * 60 * 60);

    const minutes = Math.floor(diff / (1000 * 60));

    // 返回格式化的字符串
    return `${days}天 ${hours}小时 ${minutes}分钟`;
  }

  const cpuPercent = computed(()=>{
    if (pageData.value.cpu_status){
      let total = pageData.value.cpu_status.percent.reduce((acc:number, cur:number)=>{
        return acc + cur;
      },0)
      return (total / pageData.value.cpu_status.counts).toFixed(2);
    }else{
      return 0
    }
  })

  onMounted(()=>{
    // 基于准备好的dom，初始化echarts实例
    netIoChart = echarts.init(document.getElementById('netIo'));
    diskIoChart = echarts.init(document.getElementById('diskIo'));
    // 绘制图表
    getPageData()
    // getPageData()
  })
  onBeforeUnmount(function (){
    clearInterval(jsq)
  })
</script>

<style scoped lang="less">
.sys_info{
  display: flex;
  gap: 15px 50px;
  flex-wrap: wrap;
  .system_info_item{
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 250px;
    .title{
      font-size:16px;
      text-align:center;
      margin-bottom: 10px;
      font-weight: bold;
    }
    .detail_info{
      font-size: 12px;
    }
  }
}
.ioBox{
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 3%;
  .label{
    font-size: 14px;
    color: #666;
    margin-bottom: 5px;
    .dot{
      display: inline-block;
      width: 10px;
      height: 10px;
      border-radius: 50%;
      margin-right: 5px;
    }
  }
  .value{
    font-size:14px;
  }
}
</style>
