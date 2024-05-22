#!/usr/bin/env python
# -*-coding:utf-8 -*-

import flet as ft
from flet_core import ControlEvent

from service.common import S_Text
from service.kafka_service import kafka_service


class Broker(object):
    """
    Cluster页的组件
    kafka版本、操作系统、集群信息等等
    """
    throttle_time_ms = 'throttle_time_ms'
    cluster_id = 'cluster_id'
    controller_id = 'controller_id'
    version = 'api_version'
    detail_configs = 'detail_configs'

    def __init__(self):
        self.base_info = None
        self.api_version = None
        self.meta = None
        self.cluster_table = None

        # if not kafka_service.kac:
        #     raise Exception("请先选择一个可用的kafka连接！")

        # 先加载框架
        self.base_info_tab = ft.Tab(
            text="基础信息", content=ft.Container(), icon=ft.icons.INFO_OUTLINE,
        )

        self.node_tab = ft.Tab(
            text='集群节点列表', content=ft.Container(), icon=ft.icons.HIVE_OUTLINED,
        )

        self.config_tab = ft.Tab(
            text='Broker配置', content=ft.Container(content=ft.Text("请从broker的配置按钮进入", size=20)),
            icon=ft.icons.CONSTRUCTION_OUTLINED,
        )

        self.tab = ft.Tabs(
            animation_duration=300,
            tabs=[
                self.node_tab,
                self.base_info_tab,
                self.config_tab,
            ],
            expand=1,
        )

        self.controls = [
            self.tab
        ]

    def init(self, page=None):
        if not kafka_service.kac:
            return "请先选择一个可用的kafka连接！\nPlease select an available kafka connection first!"

        self.meta, self.api_version = kafka_service.get_brokers()
        self.base_info = ft.Column(
            [
                ft.DataTable(
                    columns=[
                        ft.DataColumn(S_Text(f"{self.throttle_time_ms}")),
                        ft.DataColumn(S_Text(f"{self.cluster_id}")),
                        ft.DataColumn(S_Text(f"{self.controller_id}")),
                        ft.DataColumn(S_Text(f"{self.version}")),
                    ],
                    rows=[
                        ft.DataRow(
                            cells=[
                                ft.DataCell(S_Text(f"{self.meta[self.throttle_time_ms]}")),
                                ft.DataCell(S_Text(f"{self.meta[self.cluster_id]}")),
                                ft.DataCell(S_Text(f"{self.meta[self.controller_id]}")),
                                ft.DataCell(S_Text(f"{self.api_version}")),
                            ],
                        )
                    ],
                    column_spacing=20,

                )
            ],
            scroll=ft.ScrollMode.ALWAYS,
            height=page.window_height,

        )

        self.cluster_table = ft.Column(
            [
                ft.DataTable(
                    columns=[
                        ft.DataColumn(S_Text("broker.id")),
                        ft.DataColumn(S_Text("host")),
                        ft.DataColumn(S_Text("port")),
                        ft.DataColumn(S_Text("机架感知")),
                        ft.DataColumn(S_Text("查看配置")),
                    ],
                    rows=[
                        ft.DataRow(
                            cells=[
                                ft.DataCell(S_Text(broker['node_id'])),
                                ft.DataCell(S_Text(broker['host'])),
                                ft.DataCell(S_Text(broker['port'])),
                                ft.DataCell(S_Text(broker['rack'])),
                                ft.DataCell(ft.IconButton(icon=ft.icons.CONSTRUCTION, data=broker['node_id'],
                                                          on_click=self.show_config_tab)),
                            ],
                        )
                        for broker in sorted(self.meta['brokers'], key=lambda x: x['node_id'], )
                    ],
                    column_spacing=20,
                )
            ],
            scroll=ft.ScrollMode.ALWAYS,
            height=page.window_height,

        )

        self.base_info_tab.content = ft.Row(
            wrap=False,  # 禁止换行，以确保内容在一行内展示并出现滚动条
            scroll=ft.ScrollMode.ALWAYS,  # 设置滚动条始终显示
            expand=True,  # 让Row填充页面宽度
            controls=[
                ft.Container(
                    content=self.base_info, alignment=ft.alignment.top_left, padding=10
                )])

        self.node_tab.content = ft.Row(
            wrap=False,  # 禁止换行，以确保内容在一行内展示并出现滚动条
            scroll=ft.ScrollMode.ALWAYS,  # 设置滚动条始终显示
            expand=True,  # 让Row填充页面宽度
            controls=[ft.Container(
                self.cluster_table, alignment=ft.alignment.top_left, padding=10, adaptive=True
            )])

    def show_config_tab(self, e: ControlEvent):
        """
        打开侧边栏
        """
        e.control.disabled = True
        broker_id = e.control.data
        configs = kafka_service.get_configs(res_type='broker', name=broker_id)

        _col = ft.ListView(expand=True, spacing=10, padding=10)

        for config in configs:
            config_names = f"{config['config_names']}"
            config_value = f"{config['config_value']}" if config['config_value'] is not None else ""
            _col.controls.append(ft.Row([
                ft.Text(f" • {config_names}", style=ft.TextStyle(weight=ft.FontWeight.BOLD), ),
                ft.Text(f"    值：{config_value}"),
            ]))
        self.config_tab.content = _col

        self.tab.selected_index = 2
        e.control.disabled = False
        e.page.update()
