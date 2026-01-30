const { createApp } = Vue;

// Axios 璇锋眰鎷︽埅鍣細鑷姩娣诲姞 JWT Token
axios.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },
    error => Promise.reject(error)
);

// Axios 鍝嶅簲鎷︽埅鍣細澶勭悊 401 閿欒
axios.interceptors.response.use(
    response => response,
    error => {
        if (error.response && error.response.status === 401) {
            localStorage.removeItem('token');
        }
        return Promise.reject(error);
    }
);

createApp({
    data() {
        return {
            username: '',
            password: '',
            isLoggedIn: false,
            isLoading: false,
            error: null,
            activeMenu: 'students',
            isSuperAdmin: false, // 是否为超级管理员
            syncingRole: false, // 正在同步角色信息
            enabledPermissions: [], // 存储启用的权限action_id列表
            loadingPermissions: false, // 权限加载状态
            students: [],
            coaches: [],
            accounts: [],
            // 账号筛选条件
            accountFilters: {
                id: '',
                phone: '',
                role_id: ''
            },
            // 账号管理相关数据
            showAddAccountModal: false,
            addAccountData: {
                username: '',
                password: '',
                name: '',
                phone: '',
                role_id: ''
            },
            showEditAccountModal: false,
            editAccountData: {
                id: null,
                password: '',
                name: '',
                phone: '',
                role_id: ''
            },
            activeRoles: [], // 启用的角色列表
            allRoles: [], // 所有角色列表
            orders: [],
            attributes: [],
            classifies: [],
            brands: [],
            goods: [],
            // 学生筛选条件
            studentFilters: {
                id: '',
                name: '',
                grade: '',
                status: ''
            },
            // 教练筛选条件
            coachFilters: {
                id: '',
                name: '',
                sex: '',
                subject: '',
                status: ''
            },
            // 订单筛选条件
            orderFilters: {
                id: '',
                uid: '',
                status: ''
            },
            // 子订单筛选条件
            childOrderFilters: {
                id: '',
                parentsid: '',
                goodsid: '',
                status: ''
            },
            // 属性筛选条件
            attributeFilters: {
                id: '',
                classify: '',
                status: ''
            },
            // 类型筛选条件
            classifyFilters: {
                id: '',
                level: '',
                status: ''
            },
            // 品牌筛选条件
            brandFilters: {
                id: '',
                name: '',
                status: ''
            },
            // 商品筛选条件
            goodsFilters: {
                id: '',
                name: '',
                brandid: '',
                classifyid: '',
                status: ''
            },
            // 筛选后的学生和教练数据
            filteredStudents: [],
            filteredCoaches: [],
            filteredOrders: [],
            filteredChildOrders: [],
            filteredAttributes: [],
            filteredClassifies: [],
            filteredBrands: [],
            filteredGoods: [],
            // 子订单原始数据
            childOrders: [],
            // 分页数据
            pageSize: 10,
            studentCurrentPage: 1,
            coachCurrentPage: 1,
            orderCurrentPage: 1,
            childOrderCurrentPage: 1,
            accountCurrentPage: 1,
            attributeCurrentPage: 1,
            classifyCurrentPage: 1,
            brandCurrentPage: 1,
            goodsCurrentPage: 1,
            // 弹窗状态
            showAddStudentModal: false,
            showAddCoachModal: false,
            showEditStudentModal: false,
            showEditCoachModal: false,
            showDeleteConfirm: false,
            showAddOrderDrawer: false,
            showAddOrderGoodsModal: false,
            showEditOrderDrawer: false,
            orderDrawerReadOnly: false, // 订单抽屉是否为只读模式（详情）
            showEditOrderGoodsModal: false,
            showCancelOrderConfirm: false,
            showSubmitOrderConfirm: false,
            showAddAttributeModal: false,
            showEditAttributeModal: false,
            showEnableAttributeConfirm: false,
            showDisableAttributeConfirm: false,
            showAddClassifyModal: false,
            showEditClassifyModal: false,
            showEnableClassifyConfirm: false,
            showDisableClassifyConfirm: false,
            showAttributeValuesModal: false,
            showAddBrandModal: false,
            showEditBrandModal: false,
            showEnableBrandConfirm: false,
            showDisableBrandConfirm: false,
            showAddGoodsDrawer: false,
            showEditGoodsDrawer: false,
            showEnableGoodsConfirm: false,
            showDisableGoodsConfirm: false,
            // 活动模板弹窗状态
            showAddActivityTemplateDrawer: false,
            showEditActivityTemplateDrawer: false,
            showAddTemplateGoodsModal: false,
            showEditTemplateGoodsModal: false,
            // 活动管理弹窗状态
            showAddActivityDrawer: false,
            showEditActivityDrawer: false,
            // 活动模板数据
            activityTemplates: [],
            filteredActivityTemplates: [],
            activityTemplateFilters: {
                id: '',
                name: '',
                type: '',
                status: ''
            },
            activityTemplateCurrentPage: 1,
            addActivityTemplateData: {
                name: '',
                type: '',
                select_type: '',
                status: 0,
                classify_ids: [],
                selected_classify_id: ''
            },
            templateClassifyGoods: [],
            templateClassifyGoodsCurrentPage: 1,
            // 活动模板详情
            showActivityTemplateDetailDrawer: false,
            activityTemplateDetailData: {
                id: '',
                name: '',
                type: '',
                select_type: '',
                status: 0,
                classifies: [],
                goods: []
            },
            templateDetailGoodsCurrentPage: 1,
            editActivityTemplateData: {
                id: '',
                name: '',
                type: '',
                select_type: '',
                status: 0,
                classify_ids: []
            },
            selectedTemplateGoods: [],
            editSelectedTemplateGoods: [],
            addTemplateGoodsData: {
                goods_id: '',
                name: '',
                brand_name: '',
                classify_name: '',
                price: 0
            },
            editTemplateGoodsData: {
                goods_id: '',
                name: '',
                brand_name: '',
                classify_name: '',
                price: 0
            },
            // 活动管理数据
            activities: [],
            filteredActivities: [],
            activityFilters: {
                id: '',
                name: '',
                template_id: '',
                status: ''
            },
            activityCurrentPage: 1,
            activeActivityTemplates: [],
            addActivityData: {
                name: '',
                template_id: '',
                template_type: '',
                template_select_type: '',
                start_time: '',
                end_time: '',
                status: 0,
                details: []
            },
            editActivityData: {
                id: '',
                name: '',
                template_id: '',
                template_type: '',
                template_select_type: '',
                start_time: '',
                end_time: '',
                status: 0,
                details: []
            },
            activityTemplateGoods: [],
            editActivityTemplateGoods: [],
            activityTemplateGoodsCurrentPage: 1,
            // 活动详情
            showActivityDetailDrawer: false,
            activityDetailData: {
                id: '',
                name: '',
                template_id: '',
                template_name: '',
                template_type: '',
                template_select_type: '',
                start_time: '',
                end_time: '',
                status: 0,
                details: [],
                goods: []
            },
            activityDetailGoodsCurrentPage: 1,
            // 合同管理数据
            contracts: [],
            filteredContracts: [],
            contractFilters: {
                id: '',
                student_id: '',
                student_name: '',
                type: '',
                status: '',
                payment_status: ''
            },
            contractCurrentPage: 1,
            showAddContractDrawer: false,
            showContractDetailDrawer: false,
            showSubmitContractConfirm: false,
            showRevokeContractConfirm: false,
            showTerminateContractModal: false,
            addContractData: {
                student_id: '',
                student_name: '',
                type: '',
                signature_form: '',
                name: '',
                contract_amount: '',
                signatory: ''
            },
            contractDetailData: {
                id: '',
                student_id: '',
                student_name: '',
                type: '',
                signature_form: '',
                name: '',
                contract_amount: '',
                signatory: '',
                initiating_party: '',
                initiator: '',
                status: '',
                payment_status: '',
                create_time: ''
            },
            terminationAgreementFile: null,
            // 收款管理数据
            paymentCollections: [],
            filteredPaymentCollections: [],
            paymentCollectionTab: 'regular',
            paymentCollectionSubTab: 'received',
            paymentCollectionFilters: {
                id: '',
                student_id: '',
                order_id: '',
                payer: '',
                payment_method: '',
                trading_date: '',
                status: ''
            },
            paymentCollectionCurrentPage: 1,
            showAddPaymentCollectionModal: false,
            showConfirmPaymentModal: false,
            showDeletePaymentModal: false,
            currentPaymentCollection: null,
            studentUnpaidOrders: [],
            selectedOrderInfo: null,
            paymentCollectionForm: {
                student_id: '',
                order_id: '',
                pending_amount: '',
                payment_scenario: '',
                payment_method: '',
                payment_amount: '',
                payer: '',
                payee_entity: '',
                merchant_order: '',
                trading_hours: ''
            },
            // 待认领收款数据
            unclaimedPayments: [],
            filteredUnclaimedPayments: [],
            unclaimedFilters: {
                id: '',
                payer: '',
                payment_method: '',
                arrival_date: '',
                claimer: '',
                status: ''
            },
            unclaimedCurrentPage: 1,
            showClaimConfirmModal: false,
            showDeleteUnclaimedModal: false,
            currentUnclaimed: null,
            claimOrderId: '',
            // 淘宝收款数据
            taobaoPayments: [],
            filteredTaobaoPayments: [],
            taobaoPaymentFilters: {
                id: '',
                student_id: '',
                order_id: '',
                order_time: '',
                status: ''
            },
            taobaoPaymentCurrentPage: 1,
            showAddTaobaoPaymentModal: false,
            showConfirmTaobaoPaymentModal: false,
            showDeleteTaobaoPaymentModal: false,
            currentTaobaoPayment: null,
            taobaoPaymentForm: {
                order_id: '',
                student_id: '',
                payer: '',
                zhifubao_account: '',
                payment_amount: '',
                order_time: '',
                merchant_order: '',
                pending_amount: ''
            },
            // 淘宝待认领数据
            taobaoUnclaimedPayments: [],
            filteredTaobaoUnclaimedPayments: [],
            taobaoUnclaimedFilters: {
                id: '',
                arrival_time: '',
                status: ''
            },
            taobaoUnclaimedCurrentPage: 1,
            showClaimTaobaoModal: false,
            showDeleteTaobaoUnclaimedModal: false,
            currentTaobaoUnclaimed: null,
            claimTaobaoOrderId: '',
            taobaoSubTab: 'paid',
            // 分账明细数据
            separateAccounts: [],
            filteredSeparateAccounts: [],
            separateAccountFilters: {
                id: '',
                uid: '',
                orders_id: '',
                childorders_id: '',
                goods_id: '',
                payment_id: '',
                payment_type: '',
                type: ''
            },
            separateAccountCurrentPage: 1,
            loadingSeparateAccounts: false,
            // 退款申请弹窗
            showRefundDialog: false,
            refundForm: {
                order_id: null,
                student_id: null,
                student_name: '',
                grade: '',
                gender: '',
                child_orders: [],
                selected_refunds: [],
                payments: [],
                refund_total: 0,
                unallocated_amount: 0,
                supplement_generated: false,
                taobao_supplement: null,
                regular_supplements: []
            },
            // 退款订单管理
            refundOrders: [],
            filteredRefundOrders: [],
            refundOrderFilters: {
                id: '',
                uid: '',
                order_id: ''
            },
            refundOrderCurrentPage: 1,
            loadingRefundOrders: false,
            showRefundOrderDetail: false,
            refundOrderDetail: {
                id: '',
                order_id: '',
                student_id: '',
                student_name: '',
                grade: '',
                gender: '',
                refund_amount: '',
                submitter: '',
                submit_time: '',
                status: 0,
                refund_items: [],
                refund_payments: []
            },
            // 子退费订单数据
            refundChildOrders: [],
            filteredRefundChildOrders: [],
            refundChildOrderFilters: {
                id: '',
                student_id: '',
                order_id: '',
                goods_id: '',
                status: ''
            },
            refundChildOrderCurrentPage: 1,
            loadingRefundChildOrders: false,
            // 退费管理数据
            refundManagementTab: 'regular',
            // 常规退费数据
            regularRefunds: [],
            filteredRegularRefunds: [],
            regularRefundFilters: {
                id: '',
                student_id: '',
                refund_order_id: '',
                payer: '',
                status: ''
            },
            regularRefundCurrentPage: 1,
            loadingRegularRefunds: false,
            // 淘宝退费数据
            taobaoRefunds: [],
            filteredTaobaoRefunds: [],
            taobaoRefundFilters: {
                id: '',
                student_id: '',
                refund_order_id: '',
                status: ''
            },
            taobaoRefundCurrentPage: 1,
            loadingTaobaoRefunds: false,
            // 退费明细数据
            refundPaymentDetails: [],
            filteredRefundPaymentDetails: [],
            refundPaymentDetailFilters: {
                id: '',
                student_id: '',
                order_id: '',
                refund_order_id: '',
                payment_id: '',
                payment_type: '',
                status: ''
            },
            refundPaymentDetailCurrentPage: 1,
            loadingRefundPaymentDetails: false,
            // 审批流类型数据
            approvalFlowTypes: [],
            filteredApprovalFlowTypes: [],
            approvalFlowTypeFilters: {
                id: '',
                name: '',
                status: ''
            },
            approvalFlowTypeCurrentPage: 1,
            loadingApprovalFlowTypes: false,
            showEnableFlowTypeConfirm: false,
            showDisableFlowTypeConfirm: false,
            currentFlowType: null,
            // 审批流模板数据
            approvalFlowTemplates: [],
            filteredApprovalFlowTemplates: [],
            approvalFlowTemplateFilters: {
                id: '',
                approval_flow_type_id: '',
                name: '',
                status: ''
            },
            approvalFlowTemplateCurrentPage: 1,
            loadingApprovalFlowTemplates: false,
            showAddApprovalFlowTemplateDrawer: false,
            showEnableFlowTemplateConfirm: false,
            showDisableFlowTemplateConfirm: false,
            currentFlowTemplate: null,
            addApprovalFlowTemplateData: {
                name: '',
                approval_flow_type_id: '',
                nodes: [{
                    name: '',
                    type: 0,
                    approvers: ['']
                }],
                copy_users: []
            },
            activeApprovalFlowTypes: [],
            activeUserAccounts: [],
            showViewApprovalFlowTemplateDrawer: false,
            viewApprovalFlowTemplateData: {
                name: '',
                approval_flow_type_id: '',
                flow_type_name: '',
                nodes: [],
                copy_users: []
            },
            // 审批流管理数据
            approvalFlowManagementTab: 'initiated', // 当前激活的tab：initiated, pending, completed, copied
            // 我发起的
            initiatedFlows: [],
            initiatedFlowFilters: {
                id: '',
                approval_flow_type_id: '',
                status: ''
            },
            initiatedFlowCurrentPage: 1,
            loadingInitiatedFlows: false,
            showCancelFlowConfirm: false,
            currentFlow: null,
            // 待我审批
            pendingFlows: [],
            pendingFlowFilters: {
                id: '',
                approval_flow_id: '',
                approval_flow_type_id: ''
            },
            pendingFlowCurrentPage: 1,
            loadingPendingFlows: false,
            showApprovalDetailDrawer: false,
            currentApprovalDetail: null,
            // 处理完成
            completedFlows: [],
            completedFlowFilters: {
                id: '',
                approval_flow_id: '',
                approval_flow_type_id: ''
            },
            completedFlowCurrentPage: 1,
            loadingCompletedFlows: false,
            // 抄送我的
            copiedFlows: [],
            copiedFlowFilters: {
                id: '',
                approval_flow_id: '',
                approval_flow_type_id: ''
            },
            copiedFlowCurrentPage: 1,
            loadingCopiedFlows: false,
            // 权限管理数据
            permissions: [],
            filteredPermissions: [],
            permissionFilters: {
                id: '',
                menu_id: ''
            },
            secondLevelMenus: [],
            permissionCurrentPage: 1,
            loadingPermissions: false,
            showEnablePermissionConfirm: false,
            showDisablePermissionConfirm: false,
            currentPermission: null,
            // 新增学生数据
            addStudentData: {
                name: '',
                sex_id: '',
                phone: '',
                grade_id: '',
                coach_ids: []
            },
            // 新增教练数据
            addCoachData: {
                name: '',
                sex_id: '',
                phone: '',
                subject_id: '',
                student_ids: []
            },
            // 编辑数据
            editStudentData: {
                id: '',
                name: '',
                sex_id: '',
                phone: '',
                grade_id: ''
            },
            editCoachData: {
                id: '',
                name: '',
                sex_id: '',
                phone: '',
                subject_id: ''
            },
            // 新增订单数据
            addOrderData: {
                student_id: '',
                student_name: '',
                expected_payment_time: '',      // 新增：预计付款时间
                participating_activities: '',   // 新增：参加活动显示文本
                activity_ids: [],               // 新增：活动ID数组
                discount_amount: 0              // 新增：优惠金额
            },
            // 选择的订单商品列表
            selectedOrderGoods: [],
            // 可选商品列表（用于订单选择商品）
            availableGoodsForOrder: [],
            // 新增订单商品子弹窗数据
            addOrderGoodsData: {
                goods_id: '',
                name: '',
                brand_name: '',
                classify_name: '',
                attributes: '',
                price: 0,
                total_price: 0,
                isgroup: 1
            },
            // 编辑订单数据
            editOrderData: {
                id: '',
                student_id: '',
                student_name: '',
                expected_payment_time: '',      // 新增：预计付款时间
                participating_activities: '',   // 新增：参加活动显示文本
                activity_ids: [],               // 新增：活动ID数组
                discount_amount: 0              // 新增：优惠金额
            },
            // 编辑订单选择的商品列表
            editSelectedOrderGoods: [],
            // 编辑订单商品子弹窗数据
            editOrderGoodsData: {
                goods_id: '',
                name: '',
                brand_name: '',
                classify_name: '',
                attributes: '',
                price: 0,
                total_price: 0,
                isgroup: 1
            },
            // 新增：订单活动相关数据
            orderActivities: [],          // 可用活动列表
            orderActivityTemplateGoods: [], // 订单活动模板包含的商品列表
            childDiscounts: {},           // 子订单优惠分摊 {goods_id: amount}
            editChildDiscounts: {},       // 编辑订单的子订单优惠分摊
            // 新增属性数据
            addAttributeData: {
                name: '',
                classify: ''
            },
            // 编辑属性数据
            editAttributeData: {
                id: '',
                name: '',
                classify: ''
            },
            // 新增类型数据
            addClassifyData: {
                name: '',
                level: '',
                parent_id: ''
            },
            // 编辑类型数据
            editClassifyData: {
                id: '',
                name: '',
                level: '',
                parent_id: ''
            },
            // 父级类型列表
            parentClassifies: [],
            // 类型操作ID
            classifyId: null,
            // 新增品牌数据
            addBrandData: {
                name: ''
            },
            // 编辑品牌数据
            editBrandData: {
                id: '',
                name: ''
            },
            // 品牌操作ID
            brandId: null,
            // 新增商品数据
            addGoodsData: {
                name: '',
                brandid: '',
                classifyid: '',
                isgroup: 1,  // 0=组合售卖，1=非组合售卖，默认为非组合
                price: '',
                attributevalue_ids: []
            },
            // 编辑商品数据
            editGoodsData: {
                id: '',
                name: '',
                brandid: '',
                classifyid: '',
                isgroup: 1,
                price: '',
                attributevalue_ids: []
            },
            // 商品操作ID
            goodsId: null,
            // 启用的品牌、类型、属性列表
            activeBrands: [],
            activeClassifies: [],
            activeAttributes: [],
            // 商品属性选择数据（三列布局：属性 | 属性值 | +/-）
            goodsAttributeRows: [{ attributeId: '', valueId: '' }],
            // 商品规格选择数据（三列布局：规格 | 规格值 | +/-）
            goodsSpecRows: [{ attributeId: '', valueId: '' }],
            // 可用于组合的商品列表
            availableGoodsForCombo: [],
            // 已选择的包含商品列表（前端展示用）
            selectedIncludedGoods: [],
            // 子弹窗状态
            showAddIncludedGoodsModal: false,
            // 子弹窗数据
            addIncludedGoodsData: {
                goods_id: '',
                goods_name: '',
                brand_name: '',
                classify_name: '',
                attributes: '',
                price: ''
            },
            // 属性值数据
            currentAttributeId: null,
            currentAttributeName: '',
            attributeValues: [''],
            // 删除确认数据
            deleteId: null,
            deleteType: '',
            cancelOrderId: null,
            attributeId: null,
            // 性别列表和启用的教练列表
            sexes: [],
            activeCoaches: [],
            activeStudents: [],
            activeGrades: [],
            activeSubjects: [],
            activeStudentsForOrder: [],
            // 展开的菜单列表
            expandedMenus: [],
            // 菜单树数据
            menuTree: [],
            // 各列表的loading状态
            loadingStudents: false,
            loadingCoaches: false,
            loadingOrders: false,
            loadingChildOrders: false,
            loadingGoods: false,
            loadingAttributes: false,
            loadingBrands: false,
            loadingClassifies: false,
            loadingAccounts: false,
            loadingActivityTemplates: false,
            loadingActivities: false,
            loadingContracts: false,
            loadingPaymentCollections: false,
            loadingUnclaimedPayments: false,
            loadingTaobaoPayments: false,
            loadingTaobaoUnclaimed: false,
            // 角色管理数据
            roles: [],
            filteredRoles: [],
            roleFilters: {
                id: '',
                name: '',
                status: ''
            },
            roleCurrentPage: 1,
            loadingRoles: false,
            showAddRoleModal: false,
            showEditRoleModal: false,
            showRolePermissionsModal: false,
            showEnableRoleConfirm: false,
            showDisableRoleConfirm: false,
            currentRole: null,
            addRoleData: {
                name: '',
                comment: ''
            },
            editRoleData: {
                id: '',
                name: '',
                comment: ''
            },
            permissionsTree: [],
            selectedRolePermissions: [],
            availablePermissions: [],
            // 菜单管理数据
            menuManagementList: [],
            filteredMenuManagementList: [],
            menuManagementFilters: {
                id: '',
                name: '',
                status: ''
            },
            menuManagementCurrentPage: 1,
            loadingMenuManagement: false,
            showEditMenuModal: false,
            showEnableMenuConfirm: false,
            showDisableMenuConfirm: false,
            currentMenu: null,
            editMenuData: {
                id: '',
                name: '',
                sort_order: ''
            }
        };
    },
    computed: {
        // 获取属性列表（classify=0）
        attributesList() {
            return this.activeAttributes.filter(attr => attr.classify === 0);
        },
        // 获取规格列表（classify=1）
        specsList() {
            return this.activeAttributes.filter(attr => attr.classify === 1);
        },
        // 学生分页数据
        paginatedStudents() {
            const start = (this.studentCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredStudents.slice(start, end);
        },
        studentTotalPages() {
            return Math.ceil(this.filteredStudents.length / this.pageSize);
        },
        // 教练分页数据
        paginatedCoaches() {
            const start = (this.coachCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredCoaches.slice(start, end);
        },
        coachTotalPages() {
            return Math.ceil(this.filteredCoaches.length / this.pageSize);
        },
        // 订单分页数据
        paginatedOrders() {
            const start = (this.orderCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredOrders.slice(start, end);
        },
        orderTotalPages() {
            return Math.ceil(this.filteredOrders.length / this.pageSize);
        },
        // 子订单分页数据
        paginatedChildOrders() {
            const start = (this.childOrderCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredChildOrders.slice(start, end);
        },
        childOrderTotalPages() {
            return Math.ceil(this.filteredChildOrders.length / this.pageSize);
        },
        // 订单应收金额（所有商品的商品总价之和）
        orderTotalReceivable() {
            if (!this.selectedOrderGoods || this.selectedOrderGoods.length === 0) {
                return 0;
            }
            return this.selectedOrderGoods.reduce((sum, goods) => {
                return sum + parseFloat(goods.total_price || 0);
            }, 0);
        },
        // 订单优惠金额
        orderDiscountAmount() {
            return parseFloat(this.addOrderData.discount_amount || 0);
        },
        // 订单实收金额（所有商品的标准售价之和 - 优惠金额）
        orderTotalReceived() {
            if (!this.selectedOrderGoods || this.selectedOrderGoods.length === 0) {
                return 0;
            }
            const subtotal = this.selectedOrderGoods.reduce((sum, goods) => {
                return sum + parseFloat(goods.price || 0);
            }, 0);
            return subtotal - this.orderDiscountAmount;
        },
        // 过滤掉已选择的商品（用于订单选择商品子弹窗）
        availableGoodsForOrderSelection() {
            const selectedIds = this.selectedOrderGoods.map(g => g.goods_id);
            return this.availableGoodsForOrder.filter(g => !selectedIds.includes(g.id));
        },
        // 编辑订单应收金额
        editOrderTotalReceivable() {
            if (!this.editSelectedOrderGoods || this.editSelectedOrderGoods.length === 0) {
                return 0;
            }
            return this.editSelectedOrderGoods.reduce((sum, goods) => {
                return sum + parseFloat(goods.total_price || 0);
            }, 0);
        },
        // 编辑订单优惠金额
        editOrderDiscountAmount() {
            return parseFloat(this.editOrderData.discount_amount || 0);
        },
        // 编辑订单实收金额（所有商品的标准售价之和 - 优惠金额）
        editOrderTotalReceived() {
            if (!this.editSelectedOrderGoods || this.editSelectedOrderGoods.length === 0) {
                return 0;
            }
            const subtotal = this.editSelectedOrderGoods.reduce((sum, goods) => {
                return sum + parseFloat(goods.price || 0);
            }, 0);
            return subtotal - this.editOrderDiscountAmount;
        },
        // 过滤掉已选择的商品（用于编辑订单选择商品子弹窗）
        availableGoodsForEditOrderSelection() {
            const selectedIds = this.editSelectedOrderGoods.map(g => g.goods_id);
            return this.availableGoodsForOrder.filter(g => !selectedIds.includes(g.id));
        },
        // 账号分页数据
        paginatedAccounts() {
            const start = (this.accountCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.accounts.slice(start, end);
        },
        accountTotalPages() {
            return Math.ceil(this.accounts.length / this.pageSize);
        },
        // 属性分页数据
        paginatedAttributes() {
            const start = (this.attributeCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredAttributes.slice(start, end);
        },
        attributeTotalPages() {
            return Math.ceil(this.filteredAttributes.length / this.pageSize);
        },
        // 类型分页数据
        paginatedClassifies() {
            const start = (this.classifyCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredClassifies.slice(start, end);
        },
        classifyTotalPages() {
            return Math.ceil(this.filteredClassifies.length / this.pageSize);
        },
        // 品牌分页数据
        paginatedBrands() {
            const start = (this.brandCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredBrands.slice(start, end);
        },
        brandTotalPages() {
            return Math.ceil(this.filteredBrands.length / this.pageSize);
        },
        // 商品分页数据
        paginatedGoods() {
            const start = (this.goodsCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredGoods.slice(start, end);
        },
        goodsTotalPages() {
            return Math.ceil(this.filteredGoods.length / this.pageSize);
        },
        // 计算商品总价（仅前端显示，不保存）
        totalGoodsPrice() {
            if (!this.selectedIncludedGoods || this.selectedIncludedGoods.length === 0) {
                return 0;
            }
            return this.selectedIncludedGoods.reduce((sum, goods) => {
                return sum + parseFloat(goods.price || 0);
            }, 0).toFixed(2);
        },
        // 过滤掉已选择的商品（用于子弹窗下拉框）
        availableGoodsForSelection() {
            const selectedIds = this.selectedIncludedGoods.map(g => g.id);
            return this.availableGoodsForCombo.filter(g => !selectedIds.includes(g.id));
        },
        // 活动模板分页数据
        paginatedActivityTemplates() {
            const start = (this.activityTemplateCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredActivityTemplates.slice(start, end);
        },
        activityTemplateTotalPages() {
            return Math.ceil(this.filteredActivityTemplates.length / this.pageSize);
        },
        // 活动分页数据
        paginatedActivities() {
            const start = (this.activityCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredActivities.slice(start, end);
        },
        activityTotalPages() {
            return Math.ceil(this.filteredActivities.length / this.pageSize);
        },
        // 合同分页数据
        paginatedContracts() {
            const start = (this.contractCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredContracts.slice(start, end);
        },
        contractTotalPages() {
            return Math.ceil(this.filteredContracts.length / this.pageSize);
        },
        // 收款分页数据
        paginatedPaymentCollections() {
            const start = (this.paymentCollectionCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredPaymentCollections.slice(start, end);
        },
        paymentCollectionTotalPages() {
            return Math.ceil(this.filteredPaymentCollections.length / this.pageSize);
        },
        // 待认领分页数据
        paginatedUnclaimedPayments() {
            const start = (this.unclaimedCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredUnclaimedPayments.slice(start, end);
        },
        totalUnclaimedPages() {
            return Math.ceil(this.filteredUnclaimedPayments.length / 10);
        },
        // 淘宝已付款分页数据
        paginatedTaobaoPayments() {
            const start = (this.taobaoPaymentCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredTaobaoPayments.slice(start, end);
        },
        totalTaobaoPaymentPages() {
            return Math.ceil(this.filteredTaobaoPayments.length / 10);
        },
        // 淘宝待认领分页数据
        paginatedTaobaoUnclaimed() {
            const start = (this.taobaoUnclaimedCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredTaobaoUnclaimedPayments.slice(start, end);
        },
        totalTaobaoUnclaimedPages() {
            return Math.ceil(this.filteredTaobaoUnclaimedPayments.length / 10);
        },
        // 分账明细分页数据
        paginatedSeparateAccounts() {
            const start = (this.separateAccountCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredSeparateAccounts.slice(start, end);
        },
        totalSeparateAccountPages() {
            return Math.ceil(this.filteredSeparateAccounts.length / 10);
        },
        // 退款订单分页数据
        paginatedRefundOrders() {
            const start = (this.refundOrderCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredRefundOrders.slice(start, end);
        },
        refundOrderTotalPages() {
            return Math.ceil(this.filteredRefundOrders.length / 10);
        },
        // 子退费订单分页数据
        paginatedRefundChildOrders() {
            const start = (this.refundChildOrderCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredRefundChildOrders.slice(start, end);
        },
        refundChildOrderTotalPages() {
            return Math.ceil(this.filteredRefundChildOrders.length / 10);
        },
        // 常规退费分页数据
        paginatedRegularRefunds() {
            const start = (this.regularRefundCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredRegularRefunds.slice(start, end);
        },
        regularRefundTotalPages() {
            return Math.ceil(this.filteredRegularRefunds.length / 10);
        },
        // 淘宝退费分页数据
        paginatedTaobaoRefunds() {
            const start = (this.taobaoRefundCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredTaobaoRefunds.slice(start, end);
        },
        taobaoRefundTotalPages() {
            return Math.ceil(this.filteredTaobaoRefunds.length / 10);
        },
        // 退费明细分页数据
        paginatedRefundPaymentDetails() {
            const start = (this.refundPaymentDetailCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredRefundPaymentDetails.slice(start, end);
        },
        refundPaymentDetailTotalPages() {
            return Math.ceil(this.filteredRefundPaymentDetails.length / 10);
        },
        // 审批流类型分页数据
        paginatedApprovalFlowTypes() {
            const start = (this.approvalFlowTypeCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredApprovalFlowTypes.slice(start, end);
        },
        approvalFlowTypeTotalPages() {
            return Math.ceil(this.filteredApprovalFlowTypes.length / 10);
        },
        // 审批流模板分页数据
        paginatedApprovalFlowTemplates() {
            const start = (this.approvalFlowTemplateCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredApprovalFlowTemplates.slice(start, end);
        },
        approvalFlowTemplateTotalPages() {
            return Math.ceil(this.filteredApprovalFlowTemplates.length / 10);
        },
        // 审批流管理分页数据
        paginatedInitiatedFlows() {
            const start = (this.initiatedFlowCurrentPage - 1) * 10;
            const end = start + 10;
            return this.initiatedFlows.slice(start, end);
        },
        initiatedFlowTotalPages() {
            return Math.ceil(this.initiatedFlows.length / 10);
        },
        paginatedPendingFlows() {
            const start = (this.pendingFlowCurrentPage - 1) * 10;
            const end = start + 10;
            return this.pendingFlows.slice(start, end);
        },
        pendingFlowTotalPages() {
            return Math.ceil(this.pendingFlows.length / 10);
        },
        paginatedCompletedFlows() {
            const start = (this.completedFlowCurrentPage - 1) * 10;
            const end = start + 10;
            return this.completedFlows.slice(start, end);
        },
        completedFlowTotalPages() {
            return Math.ceil(this.completedFlows.length / 10);
        },
        paginatedCopiedFlows() {
            const start = (this.copiedFlowCurrentPage - 1) * 10;
            const end = start + 10;
            return this.copiedFlows.slice(start, end);
        },
        copiedFlowTotalPages() {
            return Math.ceil(this.copiedFlows.length / 10);
        },
        // 权限管理分页数据
        paginatedPermissions() {
            const start = (this.permissionCurrentPage - 1) * 10;
            const end = start + 10;
            return this.filteredPermissions.slice(start, end);
        },
        permissionTotalPages() {
            return Math.ceil(this.filteredPermissions.length / 10);
        },
        // 计算是否可以选择线上付款（预计付款日期未过）
        canSelectOnlinePayment() {
            if (!this.selectedOrderInfo || !this.selectedOrderInfo.expected_payment_time) {
                return true;
            }
            const today = new Date();
            today.setHours(0, 0, 0, 0);
            const expectedDate = new Date(this.selectedOrderInfo.expected_payment_time);
            expectedDate.setHours(0, 0, 0, 0);
            return expectedDate >= today;
        },
        // 过滤掉已选择的商品（用于新增活动模板）
        availableGoodsForTemplate() {
            const selectedIds = this.selectedTemplateGoods.map(g => g.goods_id);
            return this.availableGoodsForOrder.filter(g => !selectedIds.includes(g.id));
        },
        // 过滤掉已选择的商品（用于编辑活动模板）
        availableGoodsForEditTemplate() {
            const selectedIds = this.editSelectedTemplateGoods.map(g => g.goods_id);
            return this.availableGoodsForOrder.filter(g => !selectedIds.includes(g.id));
        },
        // 模板类型关联商品分页
        paginatedTemplateClassifyGoods() {
            const start = (this.templateClassifyGoodsCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.templateClassifyGoods.slice(start, end);
        },
        templateClassifyGoodsTotalPages() {
            return Math.ceil(this.templateClassifyGoods.length / this.pageSize);
        },
        // 模板详情商品分页
        paginatedTemplateDetailGoods() {
            const start = (this.templateDetailGoodsCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.activityTemplateDetailData.goods.slice(start, end);
        },
        templateDetailGoodsTotalPages() {
            return Math.ceil(this.activityTemplateDetailData.goods.length / this.pageSize);
        },
        // 活动模板商品分页（新增活动抽屉）
        paginatedActivityTemplateGoods() {
            const start = (this.activityTemplateGoodsCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.activityTemplateGoods.slice(start, end);
        },
        activityTemplateGoodsTotalPages() {
            return Math.ceil(this.activityTemplateGoods.length / this.pageSize);
        },
        // 活动详情商品分页
        paginatedActivityDetailGoods() {
            const start = (this.activityDetailGoodsCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.activityDetailData.goods.slice(start, end);
        },
        activityDetailGoodsTotalPages() {
            return Math.ceil(this.activityDetailData.goods.length / this.pageSize);
        },
        // 角色管理分页数据
        paginatedRoles() {
            const start = (this.roleCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredRoles.slice(start, end);
        },
        roleTotalPages() {
            return Math.ceil(this.filteredRoles.length / this.pageSize);
        },
        // 菜单管理分页数据
        paginatedMenuManagement() {
            const start = (this.menuManagementCurrentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.filteredMenuManagementList.slice(start, end);
        },
        menuManagementTotalPages() {
            return Math.ceil(this.filteredMenuManagementList.length / this.pageSize);
        }
    },
    mounted() {
        // 页面加载时检查用户是否已登录
        this.checkLoginStatus();
        // 加载菜单数据
        this.fetchMenus();
    },
    methods: {
        // 检查登录状态
        async checkLoginStatus() {
            try {
                const response = await axios.get('/api/profile', { withCredentials: true });
                if (response.data.username) {
                    this.isLoggedIn = true;
                    this.username = response.data.data.username;
                    this.isSuperAdmin = response.data.data.is_super_admin || false;
                    this.error = null;
                    // 加载用户权限
                    await this.fetchEnabledPermissions();
                    // 登录成功后获取学生和教练数据
                    this.fetchStudents();
                    this.fetchCoaches();
                }
            } catch (err) {
                // 未登录或其他错误，不显示错误信息
                this.isLoggedIn = false;
            }
        },
        
        // 获取菜单数据
        async fetchMenus() {
            try {
                const response = await axios.get('/api/menus');
                this.menuTree = response.data.data.menus;
            } catch (err) {
                console.error('获取菜单失败:', err);
                // 如果获取失败，使用默认菜单结构
                this.menuTree = [
                    {
                        id: 1,
                        name: '学生管理',
                        route: null,
                        children: [{ id: 5, name: '学生管理', route: 'students' }]
                    },
                    {
                        id: 2,
                        name: '教练管理',
                        route: null,
                        children: [{ id: 6, name: '教练管理', route: 'coaches' }]
                    },
                    {
                        id: 3,
                        name: '订单管理',
                        route: null,
                        children: [{ id: 7, name: '订单管理', route: 'orders' }]
                    },
                    {
                        id: 4,
                        name: '系统设置',
                        route: null,
                        children: [{ id: 8, name: '账号管理', route: 'accounts' }]
                    }
                ];
            }
        },

        // 切换菜单展开/收起
        toggleMenu(menuKey) {
            const index = this.expandedMenus.indexOf(menuKey);
            if (index > -1) {
                // 如果已展开，则收起
                this.expandedMenus.splice(index, 1);
            } else {
                // 如果未展开，则展开
                this.expandedMenus.push(menuKey);
            }
        },

        // 获取菜单的key（用于展开/收起）
        getMenuKey(menu) {
            // 使用route作为key，如果没有route则使用name
            return menu.route || menu.name.toLowerCase();
        },

        // 设置活动菜单
        async setActiveMenu(menu) {
            // 先同步角色信息，如果角色有变更会自动重新加载权限
            await this.syncRoleInfo();

            this.activeMenu = menu;
            // 切换菜单时刷新对应页面的数据
            if (menu === 'students') {
                this.fetchStudents();
            } else if (menu === 'coaches') {
                this.fetchCoaches();
            } else if (menu === 'orders') {
                this.fetchOrders();
            } else if (menu === 'childorders') {
                this.fetchChildOrders();
            } else if (menu === 'accounts') {
                this.fetchAccounts();
            } else if (menu === 'attributes') {
                this.fetchAttributes();
            } else if (menu === 'classifies') {
                this.fetchClassifies();
            } else if (menu === 'brands') {
                this.fetchBrands();
            } else if (menu === 'goods') {
                this.fetchGoods();
            } else if (menu === 'activity_template') {
                this.fetchActivityTemplates();
                this.fetchClassifies();
            } else if (menu === 'activity_management') {
                this.fetchActivities();
                this.fetchActivityTemplates();
            } else if (menu === 'contract_management') {
                this.fetchContracts();
            } else if (menu === 'payment_collection') {
                this.fetchPaymentCollections();
            } else if (menu === 'separate_account') {
                this.fetchSeparateAccounts();
            } else if (menu === 'refund_orders') {
                this.fetchRefundOrders();
            } else if (menu === 'refund_childorders') {
                this.fetchRefundChildOrders();
            } else if (menu === 'refund_management') {
                if (this.refundManagementTab === 'regular') {
                    this.fetchRegularRefunds();
                } else {
                    this.fetchTaobaoRefunds();
                }
            } else if (menu === 'refund_payment_detail') {
                this.fetchRefundPaymentDetails();
            } else if (menu === 'approval_flow_type') {
                this.fetchApprovalFlowTypes();
            } else if (menu === 'approval_flow_template') {
                this.fetchApprovalFlowTemplates();
            } else if (menu === 'approval_flow_management') {
                // 加载审批流类型列表用于筛选
                this.fetchApprovalFlowTypes();
                // 根据当前tab加载对应数据
                if (this.approvalFlowManagementTab === 'initiated') {
                    this.fetchInitiatedFlows();
                } else if (this.approvalFlowManagementTab === 'pending') {
                    this.fetchPendingFlows();
                } else if (this.approvalFlowManagementTab === 'completed') {
                    this.fetchCompletedFlows();
                } else if (this.approvalFlowManagementTab === 'copied') {
                    this.fetchCopiedFlows();
                }
            } else if (menu === 'permissions') {
                this.fetchPermissions();
                this.fetchSecondLevelMenus();
            } else if (menu === 'roles') {
                this.fetchRoles();
            } else if (menu === 'menu_management') {
                this.fetchMenuManagement();
            }
        },
        
        // 获取学生数据
        async fetchStudents() {
            this.loadingStudents = true;
            try {
                const response = await axios.get('/api/students');
                this.students = response.data.students;
                // 获取数据后自动筛选
                this.filteredStudents = this.students;
            } catch (err) {
                console.error('获取学生数据失败:', err);
                this.error = '获取学生数据失败';
            } finally {
                this.loadingStudents = false;
            }
        },
        
        // 获取教练数据
        async fetchCoaches() {
            this.loadingCoaches = true;
            try {
                const response = await axios.get('/api/coaches');
                this.coaches = response.data.coaches;
                // 获取数据后自动筛选
                this.filteredCoaches = this.coaches;
            } catch (err) {
                console.error('获取教练数据失败:', err);
                this.error = '获取教练数据失败';
            } finally {
                this.loadingCoaches = false;
            }
        },
        
        // 学生筛选功能
        searchStudents() {
            this.filteredStudents = this.students.filter(student => {
                // ID 精确匹配
                const idMatch = !this.studentFilters.id || student.id === parseInt(this.studentFilters.id);
                // 姓名模糊匹配
                const nameMatch = !this.studentFilters.name || student.student_name.toLowerCase().includes(this.studentFilters.name.toLowerCase());
                // 年级精确匹配
                const gradeMatch = !this.studentFilters.grade || student.grade === this.studentFilters.grade;
                // 状态精确匹配
                const statusMatch = this.studentFilters.status === '' || student.status === parseInt(this.studentFilters.status);
                return idMatch && nameMatch && gradeMatch && statusMatch;
            });
        },
        
        // 重置学生筛选
        resetStudentFilters() {
            this.studentFilters = {
                id: '',
                name: '',
                grade: '',
                status: ''
            };
            this.filteredStudents = this.students;
        },
        
        // 教练筛选功能
        searchCoaches() {
            this.filteredCoaches = this.coaches.filter(coach => {
                // ID 精确匹配
                const idMatch = !this.coachFilters.id || coach.id === parseInt(this.coachFilters.id);
                // 姓名模糊匹配
                const nameMatch = !this.coachFilters.name || coach.coach_name.toLowerCase().includes(this.coachFilters.name.toLowerCase());
                // 性别精确匹配
                const sexMatch = !this.coachFilters.sex || coach.sex === this.coachFilters.sex;
                // 学科精确匹配
                const subjectMatch = !this.coachFilters.subject || coach.subject === this.coachFilters.subject;
                // 状态精确匹配
                const statusMatch = this.coachFilters.status === '' || coach.status === parseInt(this.coachFilters.status);
                return idMatch && nameMatch && sexMatch && subjectMatch && statusMatch;
            });
        },
        
        // 重置教练筛选
        resetCoachFilters() {
            this.coachFilters = {
                id: '',
                name: '',
                sex: '',
                subject: '',
                status: ''
            };
            this.filteredCoaches = this.coaches;
        },
        
        // 登录方法
        async login() {
            if (!this.username || !this.password) {
                this.error = '用户名和密码不能为空';
                return;
            }
            
            try {
                this.isLoading = true;
                this.error = null;
                
                const response = await axios.post('/api/login', {
                    username: this.username,
                    password: this.password
                }, { withCredentials: true });
                
                if (response.data.data && response.data.data.token) {
                    // 瀛樺偍 JWT Token
                    localStorage.setItem('token', response.data.data.token);

                    this.isLoggedIn = true;
                    this.username = response.data.data.username;
                    this.isSuperAdmin = response.data.data.is_super_admin || false; // 存储是否为超级管理员
                    this.password = ''; // 清空密码
                    // 登录成功后重新加载菜单和权限
                    await this.fetchMenus();
                    await this.fetchEnabledPermissions();
                    // 获取学生和教练数据
                    this.fetchStudents();
                    this.fetchCoaches();
                }
            } catch (err) {
                this.error = err.response?.data?.error || '登录失败，请检查用户名和密码';
                console.error('Login error:', err);
            } finally {
                this.isLoading = false;
            }
        },
        
        // 登出方法
        async logout() {
            try {
                this.isLoading = true;
                await axios.post('/api/logout', {}, { withCredentials: true });
                // 娓呴櫎 Token
                    localStorage.removeItem('token');

                    // 娓呴櫎 Token
                    localStorage.removeItem('token');

                    this.isLoggedIn = false;
                    this.username = '';
                this.error = null;
            } catch (err) {
                this.error = '登出失败，请稍后重试';
                console.error('Logout error:', err);
            } finally {
                this.isLoading = false;
            }
        },
        
        // 学生新增功能
        async openAddStudentModal() {
            // 获取性别列表
            try {
                const response = await axios.get('/api/sexes', { withCredentials: true });
                this.sexes = response.data.data.sexes;
            } catch (err) {
                console.error('获取性别列表失败:', err);
                this.error = '获取性别列表失败';
            }

            // 获取启用的教练列表
            try {
                const response = await axios.get('/api/coaches/active', { withCredentials: true });
                this.activeCoaches = response.data.coaches;
            } catch (err) {
                console.error('获取启用教练列表失败:', err);
                this.error = '获取启用教练列表失败';
            }

            // 获取启用的年级列表
            try {
                const response = await axios.get('/api/grades/active', { withCredentials: true });
                this.activeGrades = response.data.data.grades;
            } catch (err) {
                console.error('获取启用年级列表失败:', err);
                this.error = '获取启用年级列表失败';
            }

            this.showAddStudentModal = true;
        },

        closeAddStudentModal() {
            this.showAddStudentModal = false;
            this.addStudentData = {
                name: '',
                sex_id: '',
                phone: '',
                grade_id: '',
                coach_ids: []
            };
        },

        async saveAddStudent() {
            // 验证必填字段
            if (!this.addStudentData.name || !this.addStudentData.sex_id ||
                !this.addStudentData.phone || !this.addStudentData.grade_id) {
                alert('请填写所有必填字段');
                return;
            }

            try {
                // 调用后端API来新增学生
                const response = await axios.post('/api/students', {
                    student_name: this.addStudentData.name,
                    sex_id: this.addStudentData.sex_id,
                    phone: this.addStudentData.phone,
                    grade_id: this.addStudentData.grade_id,
                    coach_ids: this.addStudentData.coach_ids
                }, { withCredentials: true });

                if (response.data.message === '学生添加成功') {
                    // 添加成功后刷新学生数据
                    await this.fetchStudents();
                    this.closeAddStudentModal();
                    alert('学生添加成功');
                }
            } catch (err) {
                console.error('新增学生失败:', err);
                alert(err.response?.data?.error || '新增学生失败');
            }
        },

        // 教练新增功能
        async openAddCoachModal() {
            // 获取性别列表
            try {
                const response = await axios.get('/api/sexes', { withCredentials: true });
                this.sexes = response.data.data.sexes;
            } catch (err) {
                console.error('获取性别列表失败:', err);
                this.error = '获取性别列表失败';
            }

            // 获取启用的学生列表
            try {
                const response = await axios.get('/api/students/active', { withCredentials: true });
                this.activeStudents = response.data.students;
            } catch (err) {
                console.error('获取启用学生列表失败:', err);
                this.error = '获取启用学生列表失败';
            }

            // 获取启用的学科列表
            try {
                const response = await axios.get('/api/subjects/active', { withCredentials: true });
                this.activeSubjects = response.data.data.subjects;
            } catch (err) {
                console.error('获取启用学科列表失败:', err);
                this.error = '获取启用学科列表失败';
            }

            this.showAddCoachModal = true;
        },

        closeAddCoachModal() {
            this.showAddCoachModal = false;
            this.addCoachData = {
                name: '',
                sex_id: '',
                phone: '',
                subject_id: '',
                student_ids: []
            };
        },

        async saveAddCoach() {
            // 验证必填字段
            if (!this.addCoachData.name || !this.addCoachData.sex_id ||
                !this.addCoachData.phone || !this.addCoachData.subject_id) {
                alert('请填写所有必填字段');
                return;
            }

            try {
                // 调用后端API来新增教练
                const response = await axios.post('/api/coaches', {
                    coach_name: this.addCoachData.name,
                    sex_id: this.addCoachData.sex_id,
                    phone: this.addCoachData.phone,
                    subject_id: this.addCoachData.subject_id,
                    student_ids: this.addCoachData.student_ids
                }, { withCredentials: true });

                if (response.data.message === '教练添加成功') {
                    // 添加成功后刷新教练数据
                    await this.fetchCoaches();
                    this.closeAddCoachModal();
                    alert('教练添加成功');
                }
            } catch (err) {
                console.error('新增教练失败:', err);
                alert(err.response?.data?.error || '新增教练失败');
            }
        },

        // 学生编辑功能
        async openEditStudentModal(student) {
            // 获取性别列表
            try {
                const response = await axios.get('/api/sexes', { withCredentials: true });
                this.sexes = response.data.data.sexes;
            } catch (err) {
                console.error('获取性别列表失败:', err);
                this.error = '获取性别列表失败';
            }

            // 获取启用的年级列表
            try {
                const response = await axios.get('/api/grades/active', { withCredentials: true });
                this.activeGrades = response.data.data.grades;
            } catch (err) {
                console.error('获取启用年级列表失败:', err);
                this.error = '获取启用年级列表失败';
            }

            this.editStudentData = {
                id: student.id,
                name: student.student_name,
                sex_id: student.sex_id,
                phone: student.phone,
                grade_id: student.grade_id
            };
            this.showEditStudentModal = true;
        },

        closeEditStudentModal() {
            this.showEditStudentModal = false;
            this.editStudentData = {
                id: '',
                name: '',
                sex_id: '',
                phone: '',
                grade_id: ''
            };
        },

        async saveEditStudent() {
            try {
                // 调用后端API来更新学生信息
                const response = await axios.put(`/api/students/${this.editStudentData.id}`, {
                    student_name: this.editStudentData.name,
                    sex_id: this.editStudentData.sex_id,
                    phone: this.editStudentData.phone,
                    grade_id: this.editStudentData.grade_id
                }, { withCredentials: true });

                if (response.data.message === '学生信息更新成功') {
                    // 更新成功后刷新学生数据
                    await this.fetchStudents();
                    this.closeEditStudentModal();
                }
            } catch (err) {
                console.error('更新学生信息失败:', err);
                this.error = err.response?.data?.error || '更新学生信息失败';
            }
        },
        
        // 教练编辑功能
        async openEditCoachModal(coach) {
            // 获取性别列表
            try {
                const response = await axios.get('/api/sexes', { withCredentials: true });
                this.sexes = response.data.data.sexes;
            } catch (err) {
                console.error('获取性别列表失败:', err);
                this.error = '获取性别列表失败';
            }

            // 获取启用的学科列表
            try {
                const response = await axios.get('/api/subjects/active', { withCredentials: true });
                this.activeSubjects = response.data.data.subjects;
            } catch (err) {
                console.error('获取启用学科列表失败:', err);
                this.error = '获取启用学科列表失败';
            }

            this.editCoachData = {
                id: coach.id,
                name: coach.coach_name,
                sex_id: coach.sex_id,
                phone: coach.phone,
                subject_id: coach.subject_id
            };
            this.showEditCoachModal = true;
        },

        closeEditCoachModal() {
            this.showEditCoachModal = false;
            this.editCoachData = {
                id: '',
                name: '',
                sex_id: '',
                phone: '',
                subject_id: ''
            };
        },

        async saveEditCoach() {
            try {
                // 调用后端API来更新教练信息
                const response = await axios.put(`/api/coaches/${this.editCoachData.id}`, {
                    coach_name: this.editCoachData.name,
                    sex_id: this.editCoachData.sex_id,
                    phone: this.editCoachData.phone,
                    subject_id: this.editCoachData.subject_id
                }, { withCredentials: true });

                if (response.data.message === '教练信息更新成功') {
                    // 更新成功后刷新教练数据
                    await this.fetchCoaches();
                    this.closeEditCoachModal();
                }
            } catch (err) {
                console.error('更新教练信息失败:', err);
                this.error = err.response?.data?.error || '更新教练信息失败';
            }
        },
        
        // 删除功能
        openDeleteConfirm(id, type) {
            this.deleteId = id;
            this.deleteType = type;
            this.showDeleteConfirm = true;
        },
        
        closeDeleteConfirm() {
            this.showDeleteConfirm = false;
            this.deleteId = null;
            this.deleteType = '';
        },
        
        async confirmDelete() {
            try {
                let response;
                // 根据删除类型调用相应的后端API
                if (this.deleteType === 'student') {
                    response = await axios.delete(`/api/students/${this.deleteId}`, { withCredentials: true });
                    if (response.data.message === '学生删除成功') {
                        await this.fetchStudents(); // 更新学生数据
                    }
                } else if (this.deleteType === 'coach') {
                    response = await axios.delete(`/api/coaches/${this.deleteId}`, { withCredentials: true });
                    if (response.data.message === '教练删除成功') {
                        await this.fetchCoaches(); // 更新教练数据
                    }
                } else if (this.deleteType === 'activityTemplate') {
                    await this.deleteActivityTemplate(this.deleteId);
                } else if (this.deleteType === 'activity') {
                    await this.deleteActivity(this.deleteId);
                }
                this.closeDeleteConfirm();
            } catch (err) {
                console.error('删除数据失败:', err);
                this.error = err.response?.data?.error || '删除数据失败';
                this.closeDeleteConfirm();
            }
        },

        // 学生状态更新
        async updateStudentStatus(id, status) {
            try {
                await axios.put(`/api/students/${id}/status`, {
                    status: status
                }, { withCredentials: true });
                const statusText = status === 0 ? '启用' : '禁用';
                alert(`学生状态已更新为${statusText}`);
                await this.fetchStudents();
            } catch (err) {
                console.error('更新学生状态失败:', err);
                alert(err.response?.data?.error || '更新学生状态失败');
            }
        },

        // 教练状态更新
        async updateCoachStatus(id, status) {
            try {
                await axios.put(`/api/coaches/${id}/status`, {
                    status: status
                }, { withCredentials: true });
                const statusText = status === 0 ? '启用' : '禁用';
                alert(`教练状态已更新为${statusText}`);
                await this.fetchCoaches();
            } catch (err) {
                console.error('更新教练状态失败:', err);
                alert(err.response?.data?.error || '更新教练状态失败');
            }
        },

        // 账号管理功能
        async fetchAccounts() {
            this.loadingAccounts = true;
            try {
                const params = new URLSearchParams();
                if (this.accountFilters.id) params.append('id', this.accountFilters.id);
                if (this.accountFilters.phone) params.append('phone', this.accountFilters.phone);
                if (this.accountFilters.role_id) params.append('role_id', this.accountFilters.role_id);

                const response = await axios.get(`/api/accounts?${params.toString()}`, { withCredentials: true });
                this.accounts = response.data.accounts;

                // 加载所有角色用于筛选
                await this.fetchAllRoles();
            } catch (err) {
                console.error('获取账号数据失败:', err);
                this.error = '获取账号数据失败';
            } finally {
                this.loadingAccounts = false;
            }
        },

        // 获取所有角色（用于筛选）
        async fetchAllRoles() {
            try {
                const response = await fetch('/api/roles', { credentials: 'include' });
                if (response.ok) {
                    const data = await response.json();
                    this.allRoles = data.roles || [];
                }
            } catch (err) {
                console.error('获取角色列表失败:', err);
            }
        },

        // 获取启用的角色（用于新增账号）
        async fetchActiveRoles() {
            try {
                const response = await fetch('/api/roles?status=0', { credentials: 'include' });
                if (response.ok) {
                    const data = await response.json();
                    this.activeRoles = data.roles || [];
                }
            } catch (err) {
                console.error('获取启用角色失败:', err);
            }
        },

        // 筛选账号
        filterAccounts() {
            this.fetchAccounts();
            this.accountCurrentPage = 1;
        },

        // 重置账号筛选
        resetAccountFilters() {
            this.accountFilters = {
                id: '',
                phone: '',
                role_id: ''
            };
            this.fetchAccounts();
            this.accountCurrentPage = 1;
        },

        // 打开新增账号弹窗
        async openAddAccountModal() {
            await this.fetchActiveRoles();
            this.addAccountData = {
                name: '',
                phone: '',
                role_id: ''
            };
            this.showAddAccountModal = true;
        },

        // 关闭新增账号弹窗
        closeAddAccountModal() {
            this.showAddAccountModal = false;
        },

        // 提交新增账号
        async submitAddAccount() {
            if (!this.addAccountData.username) {
                alert('请输入账号');
                return;
            }

            if (!this.addAccountData.password) {
                alert('请输入密码');
                return;
            }

            if (!this.addAccountData.name) {
                alert('请输入姓名');
                return;
            }

            if (!this.addAccountData.phone) {
                alert('请输入手机号');
                return;
            }

            if (!this.addAccountData.role_id) {
                alert('请选择角色');
                return;
            }

            try {
                const response = await fetch('/api/accounts', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include',
                    body: JSON.stringify(this.addAccountData)
                });

                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '创建账号失败');
                }

                alert('账号创建成功');
                this.closeAddAccountModal();
                // 重置表单
                this.addAccountData = {
                    username: '',
                    password: '',
                    name: '',
                    phone: '',
                    role_id: ''
                };
                this.fetchAccounts();
            } catch (err) {
                alert('创建账号失败：' + err.message);
            }
        },

        // 打开编辑账号弹窗
        async openEditAccountModal(account) {
            // 加载启用的角色列表
            await this.fetchActiveRoles();

            // 填充编辑表单数据
            this.editAccountData = {
                id: account.id,
                password: '', // 密码字段留空，让用户输入新密码
                name: account.name,
                phone: account.phone,
                role_id: account.role_id
            };

            this.showEditAccountModal = true;
        },

        // 关闭编辑账号弹窗
        closeEditAccountModal() {
            this.showEditAccountModal = false;
            // 重置表单
            this.editAccountData = {
                id: null,
                password: '',
                name: '',
                phone: '',
                role_id: ''
            };
        },

        // 提交编辑账号
        async submitEditAccount() {
            if (!this.editAccountData.password) {
                alert('请输入密码');
                return;
            }

            if (!this.editAccountData.name) {
                alert('请输入姓名');
                return;
            }

            if (!this.editAccountData.phone) {
                alert('请输入手机号');
                return;
            }

            if (!this.editAccountData.role_id) {
                alert('请选择角色');
                return;
            }

            try {
                const response = await fetch(`/api/accounts/${this.editAccountData.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include',
                    body: JSON.stringify({
                        password: this.editAccountData.password,
                        name: this.editAccountData.name,
                        phone: this.editAccountData.phone,
                        role_id: this.editAccountData.role_id
                    })
                });

                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '编辑账号失败');
                }

                alert('账号编辑成功');
                this.closeEditAccountModal();
                this.fetchAccounts();
            } catch (err) {
                alert('编辑账号失败：' + err.message);
            }
        },

        async enableAccount(id) {
            try {
                await axios.put(`/api/accounts/${id}/status`, {
                    status: 0
                }, { withCredentials: true });
                alert('账号已启用');
                await this.fetchAccounts();
            } catch (err) {
                console.error('启用账号失败:', err);
                alert(err.response?.data?.error || '启用账号失败');
            }
        },

        async disableAccount(id) {
            try {
                await axios.put(`/api/accounts/${id}/status`, {
                    status: 1
                }, { withCredentials: true });
                alert('账号已禁用');
                await this.fetchAccounts();
            } catch (err) {
                console.error('禁用账号失败:', err);
                alert(err.response?.data?.error || '禁用账号失败');
            }
        },

        // ==================== 订单管理功能 ====================

        // 获取订单数据
        async fetchOrders() {
            this.loadingOrders = true;
            try {
                const response = await axios.get('/api/orders', { withCredentials: true });
                this.orders = response.data.orders;
                this.filteredOrders = this.orders;
            } catch (err) {
                console.error('获取订单失败:', err);
                this.error = '获取订单失败';
            } finally {
                this.loadingOrders = false;
            }
        },

        // 订单筛选功能
        searchOrders() {
            this.filteredOrders = this.orders.filter(order => {
                const idMatch = !this.orderFilters.id || order.id === parseInt(this.orderFilters.id);
                const uidMatch = !this.orderFilters.uid || order.uid === parseInt(this.orderFilters.uid);
                const statusMatch = this.orderFilters.status === '' || order.status === parseInt(this.orderFilters.status);
                return idMatch && uidMatch && statusMatch;
            });
        },

        // 重置订单筛选
        resetOrderFilters() {
            this.orderFilters = {
                id: '',
                uid: '',
                status: ''
            };
            this.filteredOrders = this.orders;
        },

        // 获取订单状态文本
        getOrderStatusText(status) {
            const statusMap = {
                10: '草稿',
                20: '未支付',
                30: '部分支付',
                40: '已支付',
                50: '退费中',
                99: '已作废'
            };
            return statusMap[status] || '未知';
        },

        // 打开新增订单抽屉
        async openAddOrderDrawer() {
            // 获取启用的学生列表
            try {
                const response = await axios.get('/api/students/active', { withCredentials: true });
                this.activeStudentsForOrder = response.data.students;
            } catch (err) {
                console.error('获取启用学生列表失败:', err);
                this.error = '获取启用学生列表失败';
            }
            // 获取启用的商品列表
            try {
                const response = await axios.get('/api/goods/active-for-order', { withCredentials: true });
                this.availableGoodsForOrder = response.data.goods;
            } catch (err) {
                console.error('获取启用商品列表失败:', err);
                this.error = '获取启用商品列表失败';
            }
            this.showAddOrderDrawer = true;
        },

        // 关闭新增订单抽屉
        closeAddOrderDrawer() {
            this.showAddOrderDrawer = false;
            this.addOrderData = {
                student_id: '',
                student_name: '',
                expected_payment_time: '',
                participating_activities: '',
                activity_ids: [],
                discount_amount: 0
            };
            this.selectedOrderGoods = [];
            this.orderActivityTemplateGoods = [];
            this.childDiscounts = {};
        },

        // 预计付款时间变更处理（新增订单）
        async onExpectedPaymentTimeChange() {
            const paymentTime = this.addOrderData.expected_payment_time;

            if (!paymentTime) {
                // 清空活动和优惠
                this.addOrderData.participating_activities = '';
                this.addOrderData.activity_ids = [];
                this.addOrderData.discount_amount = 0;
                this.orderActivities = [];
                this.orderActivityTemplateGoods = [];
                this.childDiscounts = {};
                return;
            }

            try {
                // 查询该时间范围内的启用活动
                const response = await axios.get('/api/activities/by-date-range', {
                    params: { payment_time: paymentTime },
                    withCredentials: true
                });

                const { has_duplicate, duplicate_type, type_name, activities } = response.data;

                // 检测活动重复
                if (has_duplicate) {
                    alert(`活动信息重复，请联系运营修改！（重复类型：${type_name}）`);
                    this.addOrderData.expected_payment_time = '';
                    this.addOrderData.participating_activities = '';
                    this.addOrderData.activity_ids = [];
                    this.orderActivities = [];
                    this.orderActivityTemplateGoods = [];
                    return;
                }

                // 设置活动信息
                this.orderActivities = activities;
                this.addOrderData.activity_ids = activities.map(a => a.id);
                this.addOrderData.participating_activities = activities.map(a => a.name).join('，');

                // 加载活动模板包含的商品
                await this.loadOrderActivityTemplateGoods(activities);

                // 如果已经选择了商品，自动计算优惠
                if (this.selectedOrderGoods.length > 0) {
                    await this.calculateOrderDiscount();
                }

            } catch (err) {
                console.error('获取活动失败:', err);
                alert(err.response?.data?.error || '获取活动失败');
            }
        },

        // 计算订单优惠金额
        async calculateOrderDiscount() {
            if (!this.addOrderData.activity_ids || this.addOrderData.activity_ids.length === 0) {
                this.addOrderData.discount_amount = 0;
                this.childDiscounts = {};
                return;
            }

            if (this.selectedOrderGoods.length === 0) {
                this.addOrderData.discount_amount = 0;
                this.childDiscounts = {};
                return;
            }

            try {
                const goods_list = this.selectedOrderGoods.map(g => ({
                    goods_id: g.goods_id,
                    price: g.price,
                    total_price: g.total_price
                }));

                const response = await axios.post('/api/orders/calculate-discount', {
                    goods_list: goods_list,
                    activity_ids: this.addOrderData.activity_ids
                }, { withCredentials: true });

                this.addOrderData.discount_amount = response.data.total_discount;
                this.childDiscounts = response.data.child_discounts;

            } catch (err) {
                console.error('计算优惠失败:', err);
                this.addOrderData.discount_amount = 0;
                this.childDiscounts = {};
            }
        },

        // 编辑订单预计付款时间变化处理
        async onEditExpectedPaymentTimeChange() {
            const paymentTime = this.editOrderData.expected_payment_time;

            // 清空商品列表和活动信息
            this.editSelectedOrderGoods = [];
            this.editOrderData.participating_activities = '';
            this.editOrderData.activity_ids = [];
            this.editOrderData.discount_amount = 0;
            this.editChildDiscounts = {};

            if (!paymentTime) {
                return;
            }

            try {
                // 查询该时间范围内的启用活动
                const response = await axios.get('/api/activities/by-date-range', {
                    params: { payment_time: paymentTime },
                    withCredentials: true
                });

                const { has_duplicate, type_name, activities } = response.data;

                // 检测活动重复
                if (has_duplicate) {
                    alert(`活动信息重复，请联系运营修改！（重复类型：${type_name}）`);
                    this.editOrderData.expected_payment_time = '';
                    return;
                }

                // 设置活动信息
                this.editOrderData.activity_ids = activities.map(a => a.id);
                this.editOrderData.participating_activities = activities.map(a => a.name).join('，');

            } catch (err) {
                console.error('获取活动失败:', err);
                alert(err.response?.data?.error || '获取活动失败');
            }
        },

        // 计算编辑订单优惠金额
        async calculateEditOrderDiscount() {
            if (!this.editOrderData.activity_ids || this.editOrderData.activity_ids.length === 0) {
                this.editOrderData.discount_amount = 0;
                this.editChildDiscounts = {};
                return;
            }

            if (this.editSelectedOrderGoods.length === 0) {
                this.editOrderData.discount_amount = 0;
                this.editChildDiscounts = {};
                return;
            }

            try {
                const goods_list = this.editSelectedOrderGoods.map(g => ({
                    goods_id: g.goods_id,
                    price: g.price,
                    total_price: g.total_price
                }));

                const response = await axios.post('/api/orders/calculate-discount', {
                    goods_list: goods_list,
                    activity_ids: this.editOrderData.activity_ids
                }, { withCredentials: true });

                this.editOrderData.discount_amount = response.data.total_discount;
                this.editChildDiscounts = response.data.child_discounts;

            } catch (err) {
                console.error('计算优惠失败:', err);
                this.editOrderData.discount_amount = 0;
                this.editChildDiscounts = {};
            }
        },

        // 学生选择变化时（用于新增订单）
        onOrderStudentChange() {
            const student = this.activeStudentsForOrder.find(s => s.id === parseInt(this.addOrderData.student_id));
            if (student) {
                this.addOrderData.student_name = student.student_name;
            } else {
                this.addOrderData.student_name = '';
            }
        },

        // 打开新增订单商品子弹窗
        openAddOrderGoodsModal() {
            this.addOrderGoodsData = {
                goods_id: '',
                name: '',
                brand_name: '',
                classify_name: '',
                attributes: '',
                price: 0,
                total_price: 0,
                isgroup: 1
            };
            this.showAddOrderGoodsModal = true;
        },

        // 关闭新增订单商品子弹窗
        closeAddOrderGoodsModal() {
            this.showAddOrderGoodsModal = false;
            this.addOrderGoodsData = {
                goods_id: '',
                name: '',
                brand_name: '',
                classify_name: '',
                attributes: '',
                price: 0,
                total_price: 0,
                isgroup: 1
            };
        },

        // 订单商品选择变化时联动其他字段
        onOrderGoodsChange() {
            const goods = this.availableGoodsForOrder.find(g => g.id === parseInt(this.addOrderGoodsData.goods_id));
            if (goods) {
                this.addOrderGoodsData.name = goods.name;
                this.addOrderGoodsData.brand_name = goods.brand_name || '';
                this.addOrderGoodsData.classify_name = goods.classify_name || '';
                this.addOrderGoodsData.attributes = goods.attributes || '';
                this.addOrderGoodsData.price = goods.price;
                this.addOrderGoodsData.total_price = goods.total_price;
                this.addOrderGoodsData.isgroup = goods.isgroup;
            } else {
                this.addOrderGoodsData.name = '';
                this.addOrderGoodsData.brand_name = '';
                this.addOrderGoodsData.classify_name = '';
                this.addOrderGoodsData.attributes = '';
                this.addOrderGoodsData.price = 0;
                this.addOrderGoodsData.total_price = 0;
                this.addOrderGoodsData.isgroup = 1;
            }
        },

        // 保存选择的商品到订单列表
        async saveOrderGoods() {
            if (!this.addOrderGoodsData.goods_id) {
                alert('请选择商品');
                return;
            }
            // 添加到已选商品列表
            this.selectedOrderGoods.push({
                goods_id: parseInt(this.addOrderGoodsData.goods_id),
                name: this.addOrderGoodsData.name,
                brand_name: this.addOrderGoodsData.brand_name,
                classify_name: this.addOrderGoodsData.classify_name,
                attributes: this.addOrderGoodsData.attributes,
                price: parseFloat(this.addOrderGoodsData.price),
                total_price: parseFloat(this.addOrderGoodsData.total_price),
                isgroup: this.addOrderGoodsData.isgroup
            });
            this.closeAddOrderGoodsModal();
            // 自动重新计算优惠
            await this.calculateOrderDiscount();
        },

        // 从订单列表中删除商品
        async removeOrderGoods(index) {
            this.selectedOrderGoods.splice(index, 1);
            // 自动重新计算优惠
            await this.calculateOrderDiscount();
        },

        // 保存新增订单
        async saveAddOrder() {
            if (!this.addOrderData.student_id) {
                alert('请选择学生');
                return;
            }

            if (this.selectedOrderGoods.length === 0) {
                alert('必须至少选择一个商品');
                return;
            }

            try {
                const goods_list = this.selectedOrderGoods.map(g => ({
                    goods_id: g.goods_id,
                    price: g.price,
                    total_price: g.total_price
                }));

                const response = await axios.post('/api/orders', {
                    student_id: this.addOrderData.student_id,
                    expected_payment_time: this.addOrderData.expected_payment_time,   // 新增
                    activity_ids: this.addOrderData.activity_ids,                     // 新增
                    discount_amount: this.addOrderData.discount_amount,               // 新增
                    child_discounts: this.childDiscounts,                             // 新增
                    goods_list: goods_list
                }, { withCredentials: true });

                if (response.data.message === '订单创建成功') {
                    await this.fetchOrders();
                    this.closeAddOrderDrawer();
                    alert('保存成功！');
                }
            } catch (err) {
                console.error('创建订单失败:', err);
                alert(err.response?.data?.error || '创建订单失败');
            }
        },

        // 打开订单详情抽屉（只读）
        async openOrderDetailDrawer(order) {
            this.orderDrawerReadOnly = true;
            await this.openEditOrderDrawer(order);
        },

        // 打开编辑订单抽屉
        async openEditOrderDrawer(order) {
            this.orderDrawerReadOnly = false;
            // 获取启用的商品列表
            try {
                const response = await axios.get('/api/goods/active-for-order', { withCredentials: true });
                this.availableGoodsForOrder = response.data.goods;
            } catch (err) {
                console.error('获取启用商品列表失败:', err);
                this.error = '获取启用商品列表失败';
            }

            // 获取订单的商品列表
            try {
                const response = await axios.get(`/api/orders/${order.id}/goods`, { withCredentials: true });
                this.editSelectedOrderGoods = response.data.goods.map(g => ({
                    goods_id: g.goodsid,
                    name: g.goods_name,
                    brand_name: g.brand_name || '',
                    classify_name: g.classify_name || '',
                    attributes: g.attributes || '',
                    price: parseFloat(g.amount_received),
                    total_price: parseFloat(g.amount_receivable),
                    isgroup: g.isgroup
                }));
            } catch (err) {
                console.error('获取订单商品列表失败:', err);
                this.editSelectedOrderGoods = [];
            }

            // 格式化预计付款时间为datetime-local格式
            const formatDateForInput = (dateStr) => {
                if (!dateStr) return '';
                const date = new Date(dateStr);
                const year = date.getFullYear();
                const month = String(date.getMonth() + 1).padStart(2, '0');
                const day = String(date.getDate()).padStart(2, '0');
                const hours = String(date.getHours()).padStart(2, '0');
                const minutes = String(date.getMinutes()).padStart(2, '0');
                return `${year}-${month}-${day}T${hours}:${minutes}`;
            };

            this.editOrderData = {
                id: order.id,
                student_id: order.uid,
                student_name: order.student_name,
                expected_payment_time: formatDateForInput(order.expected_payment_time),
                participating_activities: '',
                activity_ids: [],
                discount_amount: parseFloat(order.discount_amount) || 0,
                status: order.status
            };

            // 如果有预计付款时间，加载活动信息（仅显示，不重新计算优惠）
            if (this.editOrderData.expected_payment_time) {
                try {
                    const response = await axios.get('/api/activities/by-date-range', {
                        params: { payment_time: this.editOrderData.expected_payment_time },
                        withCredentials: true
                    });
                    const { activities } = response.data;
                    this.editOrderData.activity_ids = activities.map(a => a.id);
                    this.editOrderData.participating_activities = activities.map(a => a.name).join('，');
                } catch (err) {
                    console.error('获取活动失败:', err);
                }
            }

            // 注意：打开编辑时不重新计算优惠，保持存储的值
            // 优惠只在用户修改预计付款时间或添加/删除商品时才重新计算

            this.showEditOrderDrawer = true;
        },

        // 关闭编辑订单抽屉
        closeEditOrderDrawer() {
            this.showEditOrderDrawer = false;
            this.editOrderData = {
                id: '',
                student_id: '',
                student_name: '',
                expected_payment_time: '',
                participating_activities: '',
                activity_ids: [],
                discount_amount: 0,
                status: ''
            };
            this.editSelectedOrderGoods = [];
            this.editChildDiscounts = {};
        },

        // 打开编辑订单商品子弹窗
        openEditOrderGoodsModal() {
            this.editOrderGoodsData = {
                goods_id: '',
                name: '',
                brand_name: '',
                classify_name: '',
                attributes: '',
                price: 0,
                total_price: 0,
                isgroup: 1
            };
            this.showEditOrderGoodsModal = true;
        },

        // 关闭编辑订单商品子弹窗
        closeEditOrderGoodsModal() {
            this.showEditOrderGoodsModal = false;
            this.editOrderGoodsData = {
                goods_id: '',
                name: '',
                brand_name: '',
                classify_name: '',
                attributes: '',
                price: 0,
                total_price: 0,
                isgroup: 1
            };
        },

        // 编辑订单商品选择变化时联动其他字段
        onEditOrderGoodsChange() {
            const goods = this.availableGoodsForOrder.find(g => g.id === parseInt(this.editOrderGoodsData.goods_id));
            if (goods) {
                this.editOrderGoodsData.name = goods.name;
                this.editOrderGoodsData.brand_name = goods.brand_name || '';
                this.editOrderGoodsData.classify_name = goods.classify_name || '';
                this.editOrderGoodsData.attributes = goods.attributes || '';
                this.editOrderGoodsData.price = goods.price;
                this.editOrderGoodsData.total_price = goods.total_price;
                this.editOrderGoodsData.isgroup = goods.isgroup;
            } else {
                this.editOrderGoodsData.name = '';
                this.editOrderGoodsData.brand_name = '';
                this.editOrderGoodsData.classify_name = '';
                this.editOrderGoodsData.attributes = '';
                this.editOrderGoodsData.price = 0;
                this.editOrderGoodsData.total_price = 0;
                this.editOrderGoodsData.isgroup = 1;
            }
        },

        // 保存选择的商品到编辑订单列表
        async saveEditOrderGoods() {
            if (!this.editOrderGoodsData.goods_id) {
                alert('请选择商品');
                return;
            }
            this.editSelectedOrderGoods.push({
                goods_id: parseInt(this.editOrderGoodsData.goods_id),
                name: this.editOrderGoodsData.name,
                brand_name: this.editOrderGoodsData.brand_name,
                classify_name: this.editOrderGoodsData.classify_name,
                attributes: this.editOrderGoodsData.attributes,
                price: parseFloat(this.editOrderGoodsData.price),
                total_price: parseFloat(this.editOrderGoodsData.total_price),
                isgroup: this.editOrderGoodsData.isgroup
            });
            this.closeEditOrderGoodsModal();
            await this.calculateEditOrderDiscount();
        },

        // 从编辑订单列表中删除商品
        async removeEditOrderGoods(index) {
            this.editSelectedOrderGoods.splice(index, 1);
            await this.calculateEditOrderDiscount();
        },

        // 保存编辑订单
        async saveEditOrder() {
            if (this.editSelectedOrderGoods.length === 0) {
                alert('必须至少选择一个商品');
                return;
            }

            try {
                const goods_list = this.editSelectedOrderGoods.map(g => ({
                    goods_id: g.goods_id,
                    price: g.price,
                    total_price: g.total_price
                }));

                const response = await axios.put(`/api/orders/${this.editOrderData.id}`, {
                    goods_list: goods_list,
                    expected_payment_time: this.editOrderData.expected_payment_time,
                    activity_ids: this.editOrderData.activity_ids,
                    discount_amount: this.editOrderData.discount_amount,
                    child_discounts: this.editChildDiscounts
                }, { withCredentials: true });

                if (response.data.message === '订单更新成功') {
                    await this.fetchOrders();
                    this.closeEditOrderDrawer();
                    alert('保存成功！');
                }
            } catch (err) {
                console.error('更新订单失败:', err);
                alert(err.response?.data?.error || '更新订单失败');
            }
        },

        // 打开作废订单确认弹窗
        openCancelOrderConfirm(orderId) {
            this.cancelOrderId = orderId;
            this.showCancelOrderConfirm = true;
        },

        // 关闭作废订单确认弹窗
        closeCancelOrderConfirm() {
            this.showCancelOrderConfirm = false;
            this.cancelOrderId = null;
        },

        // 确认作废订单
        async confirmCancelOrder() {
            try {
                const response = await axios.put(`/api/orders/${this.cancelOrderId}/cancel`, {}, { withCredentials: true });

                if (response.data.message === '订单已作废') {
                    await this.fetchOrders();
                    this.closeCancelOrderConfirm();
                    alert('订单已作废');
                }
            } catch (err) {
                console.error('作废订单失败:', err);
                alert(err.response?.data?.error || '作废订单失败');
            }
        },

        // 打开提交订单确认弹窗
        confirmSubmitOrder() {
            this.showSubmitOrderConfirm = true;
        },

        // 提交订单
        async submitOrder() {
            try {
                const response = await axios.put(`/api/orders/${this.editOrderData.id}/submit`, {}, { withCredentials: true });

                if (response.data.message === '订单已提交') {
                    await this.fetchOrders();
                    this.showSubmitOrderConfirm = false;
                    this.closeEditOrderDrawer();
                    alert('订单已提交');
                }
            } catch (err) {
                console.error('提交订单失败:', err);
                alert(err.response?.data?.error || '提交订单失败');
            }
        },

        // ==================== 退款申请功能 ====================

        // 打开退款申请弹窗
        async openRefundDialog(orderId) {
            try {
                const response = await axios.get(`/api/orders/${orderId}/refund-info`, { withCredentials: true });
                const data = response.data;

                // 初始化退款表单数据（不再初始化payments）
                this.refundForm = {
                    order_id: orderId,
                    student_id: data.order.student_id,
                    student_name: data.order.student_name,
                    grade: data.order.grade,
                    gender: data.order.gender,
                    child_orders: data.child_orders || [],
                    selected_refunds: [],
                    payments: [],  // 不再立即获取收款列表
                    refund_total: 0,
                    unallocated_amount: 0,
                    supplement_generated: false,
                    taobao_supplement: null,
                    regular_supplements: [],
                    childorder_separate_amounts: {}  // 新增：子订单分账金额
                };

                this.showRefundDialog = true;
            } catch (err) {
                console.error('获取退款信息失败:', err);
                alert(err.response?.data?.error || '获取退款信息失败');
            }
        },

        // 关闭退款申请弹窗
        closeRefundDialog() {
            this.showRefundDialog = false;
            this.refundForm = {
                order_id: null,
                student_id: null,
                student_name: '',
                grade: '',
                gender: '',
                child_orders: [],
                selected_refunds: [],
                payments: [],
                refund_total: 0,
                unallocated_amount: 0,
                supplement_generated: false,
                taobao_supplement: null,
                regular_supplements: [],
                childorder_separate_amounts: {}
            };
        },

        // 选择子订单添加到待退费区
        selectChildOrderForRefund(item) {
            // 检查是否已经添加
            const exists = this.refundForm.selected_refunds.find(r => r.childorder_id === item.childorder_id);
            if (exists) {
                alert('该子订单已添加到待退费区');
                return;
            }

            // 添加到待退费区，初始退费金额为0
            this.refundForm.selected_refunds.push({
                childorder_id: item.childorder_id,
                goods_id: item.goods_id,
                goods_name: item.goods_name,
                available_refund: item.available_refund,
                refund_amount: 0
            });
        },

        // 从待退费区删除
        removeRefundItem(index) {
            this.refundForm.selected_refunds.splice(index, 1);
            // 重新计算退费总额
            this.applyRefundTotal();
        },

        // 验证退费金额
        validateRefundAmount(item) {
            if (item.refund_amount > item.available_refund) {
                item.refund_amount = item.available_refund;
                alert('退费金额不能大于可退金额');
            }
            if (item.refund_amount < 0) {
                item.refund_amount = 0;
            }
        },

        // 应用退费总额
        async applyRefundTotal() {
            // 计算退费总额
            this.refundForm.refund_total = this.refundForm.selected_refunds.reduce((sum, item) => {
                return sum + (parseFloat(item.refund_amount) || 0);
            }, 0);

            // 调用后端API获取收款列表和分账金额
            try {
                const refund_items = this.refundForm.selected_refunds.map(item => ({
                    childorder_id: item.childorder_id,
                    refund_amount: item.refund_amount
                }));

                const response = await axios.post(
                    `/api/orders/${this.refundForm.order_id}/refund-payments`,
                    { refund_items },
                    { withCredentials: true }
                );

                // 更新收款列表
                this.refundForm.payments = (response.data.payments || []).map(p => ({
                    ...p,
                    refund_amount: 0
                }));

                // 保存分账金额数据
                this.refundForm.childorder_separate_amounts = response.data.childorder_separate_amounts || {};

                // 更新待分配金额
                this.calculateUnallocatedAmount();
            } catch (err) {
                console.error('获取收款列表失败:', err);
                alert(err.response?.data?.error || '获取收款列表失败');
            }
        },

        // 计算待分配金额
        calculateUnallocatedAmount() {
            const totalPaymentRefund = this.refundForm.payments.reduce((sum, payment) => {
                return sum + (parseFloat(payment.refund_amount) || 0);
            }, 0);

            this.refundForm.unallocated_amount = this.refundForm.refund_total - totalPaymentRefund;
        },

        // 验证收款退费金额
        validatePaymentRefundAmount(payment) {
            if (payment.refund_amount > payment.available_refund) {
                payment.refund_amount = payment.available_refund;
                alert('退费金额不能大于可退金额');
            }
            if (payment.refund_amount < 0) {
                payment.refund_amount = 0;
            }
            this.calculateUnallocatedAmount();
        },

        // 应用退费信息补充
        applyRefundSupplement() {
            // 校验待分配金额是否为0
            if (Math.abs(this.refundForm.unallocated_amount) > 0.01) {
                alert('请重新调整退费路径');
                return;
            }

            // 重置退费信息补充
            this.refundForm.taobao_supplement = null;
            this.refundForm.regular_supplements = [];

            // 筛选有退费金额的收款
            const paymentsWithRefund = this.refundForm.payments.filter(p => p.refund_amount > 0);

            // 1. 处理淘宝退费
            const taobaoPayments = paymentsWithRefund.filter(p => p.payment_type === 1);
            if (taobaoPayments.length > 0) {
                const totalTaobaoRefund = taobaoPayments.reduce((sum, p) => sum + parseFloat(p.refund_amount), 0);
                this.refundForm.taobao_supplement = {
                    alipay_account: '',
                    alipay_name: '',
                    refund_amount: totalTaobaoRefund
                };
            }

            // 2. 处理常规退费
            const regularPayments = paymentsWithRefund.filter(p => p.payment_type === 0);
            if (regularPayments.length > 0) {
                const results = [];

                // 先检查所有收款是否都有payee_entity
                console.log('=== 常规收款数据 ===');
                regularPayments.forEach(p => {
                    console.log(`收款ID=${p.payment_id}, payee_entity=${p.payee_entity}, payer=${p.payer}, is_corporate_transfer=${p.is_corporate_transfer}, refund_amount=${p.refund_amount}`);
                });

                // 按收款主体分组（必须先按主体分组）
                const entityGroups = {};
                for (const payment of regularPayments) {
                    const entity = payment.payee_entity; // 0=北京, 1=西安

                    // 检查payee_entity是否有效
                    if (entity === null || entity === undefined) {
                        console.error(`收款ID=${payment.payment_id}的payee_entity为空！`);
                        alert(`收款ID=${payment.payment_id}的收款主体为空，无法生成退费补充信息`);
                        return;
                    }

                    if (!entityGroups[entity]) {
                        entityGroups[entity] = [];
                    }
                    entityGroups[entity].push(payment);
                }

                console.log('=== 按主体分组结果 ===');
                for (const entity in entityGroups) {
                    console.log(`主体${entity}: ${entityGroups[entity].length}条收款`);
                }

                // 处理每个主体
                for (const entity in entityGroups) {
                    const entityPayments = entityGroups[entity];

                    // 分离对公转账和非对公转账
                    const corporatePayments = entityPayments.filter(p => p.is_corporate_transfer === 1);
                    const nonCorporatePayments = entityPayments.filter(p => p.is_corporate_transfer === 0);

                    // 2.1 处理对公转账 - 按付款方分组
                    const corporateGroups = {};
                    for (const payment of corporatePayments) {
                        const payer = payment.payer || '';
                        if (!corporateGroups[payer]) {
                            corporateGroups[payer] = [];
                        }
                        corporateGroups[payer].push(payment);
                    }

                    for (const payer in corporateGroups) {
                        const group = corporateGroups[payer];
                        const totalRefund = group.reduce((sum, p) => sum + parseFloat(p.refund_amount), 0);
                        results.push({
                            payee_entity: parseInt(entity),
                            is_corporate_transfer: 1,
                            payer: payer || null,
                            payer_readonly: true,
                            bank_account: '',
                            refund_amount: totalRefund
                        });
                    }

                    // 2.2 处理非对公转账
                    const nonCorporateWithPayer = nonCorporatePayments.filter(p => p.payer);
                    const nonCorporateWithoutPayer = nonCorporatePayments.filter(p => !p.payer);

                    // 2.2.1 处理有付款方的 - 按付款方分组
                    const nonCorpPayerGroups = {};
                    for (const payment of nonCorporateWithPayer) {
                        const payer = payment.payer;
                        if (!nonCorpPayerGroups[payer]) {
                            nonCorpPayerGroups[payer] = [];
                        }
                        nonCorpPayerGroups[payer].push(payment);
                    }

                    for (const payer in nonCorpPayerGroups) {
                        const group = nonCorpPayerGroups[payer];
                        const totalRefund = group.reduce((sum, p) => sum + parseFloat(p.refund_amount), 0);
                        results.push({
                            payee_entity: parseInt(entity),
                            is_corporate_transfer: 0,
                            payer: payer,
                            payer_readonly: true,
                            bank_account: '',
                            refund_amount: totalRefund
                        });
                    }

                    // 2.2.2 处理无付款方的
                    if (nonCorporateWithoutPayer.length > 0) {
                        const emptyPayerTotal = nonCorporateWithoutPayer.reduce((sum, p) => sum + parseFloat(p.refund_amount), 0);

                        // 检查同主体下是否有有付款方的非对公转账
                        const hasNonCorpWithPayer = nonCorporateWithPayer.length > 0;

                        if (hasNonCorpWithPayer) {
                            // 合并到该主体下第一条有付款方的非对公转账记录
                            const targetRecord = results.find(r =>
                                r.payee_entity === parseInt(entity) &&
                                r.is_corporate_transfer === 0 &&
                                r.payer
                            );
                            if (targetRecord) {
                                targetRecord.refund_amount += emptyPayerTotal;
                            }
                        } else {
                            // 独立成组，付款方和银行账户都可编辑
                            results.push({
                                payee_entity: parseInt(entity),
                                is_corporate_transfer: 0,
                                payer: '',
                                payer_readonly: false,
                                bank_account: '',
                                refund_amount: emptyPayerTotal
                            });
                        }
                    }
                }

                this.refundForm.regular_supplements = results;
            }

            // 标记为已生成
            this.refundForm.supplement_generated = true;
        },

        // 提交退款申请
        async submitRefundApplication() {
            // 校验
            if (this.refundForm.selected_refunds.length === 0) {
                alert('请选择需要退费的子订单');
                return;
            }

            // 检查所有退费金额是否大于0
            for (const item of this.refundForm.selected_refunds) {
                if (!item.refund_amount || item.refund_amount <= 0) {
                    alert('所有退费金额必须大于0');
                    return;
                }
            }

            // 检查待分配金额是否为0
            if (Math.abs(this.refundForm.unallocated_amount) > 0.01) {
                alert('待分配金额必须为0，请调整收款列表中的退费金额');
                return;
            }

            // 检查退费金额是否超出分账金额
            let hasExceededSeparateAmount = false;
            for (const item of this.refundForm.selected_refunds) {
                const separateAmount = this.refundForm.childorder_separate_amounts[item.childorder_id] || 0;
                if (parseFloat(item.refund_amount) > separateAmount) {
                    hasExceededSeparateAmount = true;
                    break;
                }
            }

            // 如果超出分账金额，显示特殊提示
            if (hasExceededSeparateAmount) {
                const confirmed = confirm('注意：本次提交会重新分配购买顺序！\n\n确认要提交吗？');
                if (!confirmed) {
                    return;
                }
            }

            // 如果还没有生成退费信息补充，自动生成
            if (!this.refundForm.supplement_generated) {
                this.generateRefundSupplement();
            }

            // 校验退费信息补充（如果已生成）
            if (this.refundForm.supplement_generated) {
                // 校验淘宝退费
                if (this.refundForm.taobao_supplement) {
                    if (!this.refundForm.taobao_supplement.alipay_account) {
                        alert('支付宝账号不能为空');
                        return;
                    }
                    if (!this.refundForm.taobao_supplement.alipay_name) {
                        alert('支付宝名称不能为空');
                        return;
                    }
                }

                // 校验常规退费
                if (this.refundForm.regular_supplements && this.refundForm.regular_supplements.length > 0) {
                    for (const item of this.refundForm.regular_supplements) {
                        // 如果付款方可编辑，则必须填写
                        if (!item.payer_readonly && !item.payer) {
                            alert('付款方不能为空');
                            return;
                        }
                    }
                }
            }

            try {
                const requestData = {
                    order_id: this.refundForm.order_id,
                    refund_items: this.refundForm.selected_refunds.map(item => ({
                        childorder_id: item.childorder_id,
                        goods_id: item.goods_id,
                        goods_name: item.goods_name,
                        refund_amount: item.refund_amount.toFixed(2)
                    })),
                    refund_payments: this.refundForm.payments
                        .filter(p => p.refund_amount > 0)
                        .map(p => ({
                            payment_id: p.payment_id,
                            payment_type: p.payment_type,
                            refund_amount: p.refund_amount.toFixed(2)
                        }))
                };

                // 添加退费信息补充
                if (this.refundForm.supplement_generated) {
                    if (this.refundForm.taobao_supplement) {
                        requestData.taobao_supplement = this.refundForm.taobao_supplement;
                    }
                    if (this.refundForm.regular_supplements && this.refundForm.regular_supplements.length > 0) {
                        requestData.regular_supplements = this.refundForm.regular_supplements;
                    }
                }

                const response = await axios.post('/api/refund-orders', requestData, { withCredentials: true });

                alert('退款申请提交成功');
                this.closeRefundDialog();
                await this.fetchOrders();
            } catch (err) {
                console.error('提交退款申请失败:', err);
                alert(err.response?.data?.error || '提交退款申请失败');
            }
        },

        // ==================== 退款订单管理功能 ====================

        // 获取退款订单列表
        async fetchRefundOrders() {
            this.loadingRefundOrders = true;
            try {
                const params = new URLSearchParams();
                if (this.refundOrderFilters.id) params.append('id', this.refundOrderFilters.id);
                if (this.refundOrderFilters.uid) params.append('uid', this.refundOrderFilters.uid);
                if (this.refundOrderFilters.order_id) params.append('order_id', this.refundOrderFilters.order_id);

                const response = await axios.get(`/api/refund-orders?${params.toString()}`, { withCredentials: true });
                this.refundOrders = response.data.refund_orders || [];
                this.filteredRefundOrders = this.refundOrders;
            } catch (err) {
                console.error('获取退款订单失败:', err);
                this.error = '获取退款订单失败';
            } finally {
                this.loadingRefundOrders = false;
            }
        },

        // 搜索退款订单
        searchRefundOrders() {
            this.fetchRefundOrders();
            this.refundOrderCurrentPage = 1;
        },

        // 重置退款订单筛选
        resetRefundOrderFilters() {
            this.refundOrderFilters = {
                id: '',
                uid: '',
                order_id: ''
            };
            this.fetchRefundOrders();
        },

        // 打开退款订单详情
        async openRefundOrderDetail(refundOrderId) {
            try {
                const response = await axios.get(`/api/refund-orders/${refundOrderId}`, { withCredentials: true });
                const data = response.data;

                this.refundOrderDetail = {
                    id: data.refund_order.id,
                    order_id: data.refund_order.order_id,
                    student_id: data.refund_order.student_id,
                    student_name: data.refund_order.student_name,
                    grade: data.refund_order.grade,
                    gender: data.refund_order.gender,
                    refund_amount: data.refund_order.refund_amount,
                    submitter: data.refund_order.submitter,
                    submit_time: data.refund_order.submit_time,
                    status: data.refund_order.status,
                    refund_items: data.refund_items || [],
                    refund_payments: data.refund_payments || [],
                    taobao_supplement: data.taobao_supplement || null,
                    regular_supplements: data.regular_supplements || []
                };

                this.showRefundOrderDetail = true;
            } catch (err) {
                console.error('获取退款订单详情失败:', err);
                alert(err.response?.data?.error || '获取退款订单详情失败');
            }
        },

        // 关闭退款订单详情
        closeRefundOrderDetail() {
            this.showRefundOrderDetail = false;
            this.refundOrderDetail = {
                id: '',
                order_id: '',
                student_id: '',
                student_name: '',
                grade: '',
                gender: '',
                refund_amount: '',
                submitter: '',
                submit_time: '',
                status: 0,
                refund_items: [],
                refund_payments: []
            };
        },

        // 退款订单分页
        changeRefundOrderPage(page) {
            if (page >= 1 && page <= this.refundOrderTotalPages) {
                this.refundOrderCurrentPage = page;
            }
        },

        // ==================== 子退费订单管理 ====================

        // 获取子退费订单列表
        async fetchRefundChildOrders() {
            this.loadingRefundChildOrders = true;
            try {
                const params = new URLSearchParams();
                if (this.refundChildOrderFilters.id) params.append('id', this.refundChildOrderFilters.id);
                if (this.refundChildOrderFilters.student_id) params.append('student_id', this.refundChildOrderFilters.student_id);
                if (this.refundChildOrderFilters.order_id) params.append('order_id', this.refundChildOrderFilters.order_id);
                if (this.refundChildOrderFilters.goods_id) params.append('goods_id', this.refundChildOrderFilters.goods_id);
                if (this.refundChildOrderFilters.status !== '') params.append('status', this.refundChildOrderFilters.status);

                const response = await axios.get(`/api/refund-childorders?${params.toString()}`, { withCredentials: true });
                this.refundChildOrders = response.data.refund_childorders || [];
                this.filteredRefundChildOrders = this.refundChildOrders;
            } catch (err) {
                console.error('获取子退费订单列表失败:', err);
                alert(err.response?.data?.error || '获取列表失败');
            } finally {
                this.loadingRefundChildOrders = false;
            }
        },

        // 搜索子退费订单
        searchRefundChildOrders() {
            this.refundChildOrderCurrentPage = 1;
            this.fetchRefundChildOrders();
        },

        // 重置子退费订单筛选条件
        resetRefundChildOrderFilters() {
            this.refundChildOrderFilters = {
                id: '',
                student_id: '',
                order_id: '',
                goods_id: '',
                status: ''
            };
            this.fetchRefundChildOrders();
        },

        // 更改子退费订单页码
        changeRefundChildOrderPage(page) {
            if (page >= 1 && page <= this.refundChildOrderTotalPages) {
                this.refundChildOrderCurrentPage = page;
            }
        },

        // 获取退款订单状态文本
        getRefundOrderStatusText(status) {
            const statusMap = {
                0: '待审批',
                10: '已通过',
                20: '已驳回'
            };
            return statusMap[status] || '未知';
        },

        // ==================== 退费管理 ====================

        // 获取常规退费列表
        async fetchRegularRefunds() {
            this.loadingRegularRefunds = true;
            try {
                const params = new URLSearchParams();
                if (this.regularRefundFilters.id) params.append('id', this.regularRefundFilters.id);
                if (this.regularRefundFilters.student_id) params.append('student_id', this.regularRefundFilters.student_id);
                if (this.regularRefundFilters.refund_order_id) params.append('refund_order_id', this.regularRefundFilters.refund_order_id);
                if (this.regularRefundFilters.payer) params.append('payer', this.regularRefundFilters.payer);
                if (this.regularRefundFilters.status !== '') params.append('status', this.regularRefundFilters.status);

                const response = await axios.get(`/api/refund-regular-supplements?${params.toString()}`, { withCredentials: true });
                this.regularRefunds = response.data.regular_supplements || [];
                this.filteredRegularRefunds = this.regularRefunds;
            } catch (err) {
                console.error('获取常规退费列表失败:', err);
                alert(err.response?.data?.error || '获取列表失败');
            } finally {
                this.loadingRegularRefunds = false;
            }
        },

        // 搜索常规退费
        searchRegularRefunds() {
            this.regularRefundCurrentPage = 1;
            this.fetchRegularRefunds();
        },

        // 重置常规退费筛选条件
        resetRegularRefundFilters() {
            this.regularRefundFilters = {
                id: '',
                student_id: '',
                refund_order_id: '',
                payer: '',
                status: ''
            };
            this.fetchRegularRefunds();
        },

        // 更改常规退费页码
        changeRegularRefundPage(page) {
            if (page >= 1 && page <= this.regularRefundTotalPages) {
                this.regularRefundCurrentPage = page;
            }
        },

        // 获取淘宝退费列表
        async fetchTaobaoRefunds() {
            this.loadingTaobaoRefunds = true;
            try {
                const params = new URLSearchParams();
                if (this.taobaoRefundFilters.id) params.append('id', this.taobaoRefundFilters.id);
                if (this.taobaoRefundFilters.student_id) params.append('student_id', this.taobaoRefundFilters.student_id);
                if (this.taobaoRefundFilters.refund_order_id) params.append('refund_order_id', this.taobaoRefundFilters.refund_order_id);
                if (this.taobaoRefundFilters.status !== '') params.append('status', this.taobaoRefundFilters.status);

                const response = await axios.get(`/api/refund-taobao-supplements?${params.toString()}`, { withCredentials: true });
                this.taobaoRefunds = response.data.taobao_supplements || [];
                this.filteredTaobaoRefunds = this.taobaoRefunds;
            } catch (err) {
                console.error('获取淘宝退费列表失败:', err);
                alert(err.response?.data?.error || '获取列表失败');
            } finally {
                this.loadingTaobaoRefunds = false;
            }
        },

        // 搜索淘宝退费
        searchTaobaoRefunds() {
            this.taobaoRefundCurrentPage = 1;
            this.fetchTaobaoRefunds();
        },

        // 重置淘宝退费筛选条件
        resetTaobaoRefundFilters() {
            this.taobaoRefundFilters = {
                id: '',
                student_id: '',
                refund_order_id: '',
                status: ''
            };
            this.fetchTaobaoRefunds();
        },

        // 更改淘宝退费页码
        changeTaobaoRefundPage(page) {
            if (page >= 1 && page <= this.taobaoRefundTotalPages) {
                this.taobaoRefundCurrentPage = page;
            }
        },

        // ==================== 退费明细管理 ====================

        // 获取退费明细列表
        async fetchRefundPaymentDetails() {
            this.loadingRefundPaymentDetails = true;
            try {
                const params = new URLSearchParams();
                if (this.refundPaymentDetailFilters.id) params.append('id', this.refundPaymentDetailFilters.id);
                if (this.refundPaymentDetailFilters.student_id) params.append('student_id', this.refundPaymentDetailFilters.student_id);
                if (this.refundPaymentDetailFilters.order_id) params.append('order_id', this.refundPaymentDetailFilters.order_id);
                if (this.refundPaymentDetailFilters.refund_order_id) params.append('refund_order_id', this.refundPaymentDetailFilters.refund_order_id);
                if (this.refundPaymentDetailFilters.payment_id) params.append('payment_id', this.refundPaymentDetailFilters.payment_id);
                if (this.refundPaymentDetailFilters.payment_type !== '') params.append('payment_type', this.refundPaymentDetailFilters.payment_type);
                if (this.refundPaymentDetailFilters.status !== '') params.append('status', this.refundPaymentDetailFilters.status);

                const response = await axios.get(`/api/refund-payment-details?${params.toString()}`, { withCredentials: true });
                this.refundPaymentDetails = response.data.refund_payment_details || [];
                this.filteredRefundPaymentDetails = this.refundPaymentDetails;
            } catch (err) {
                console.error('获取退费明细列表失败:', err);
                alert(err.response?.data?.error || '获取列表失败');
            } finally {
                this.loadingRefundPaymentDetails = false;
            }
        },

        // 搜索退费明细
        searchRefundPaymentDetails() {
            this.refundPaymentDetailCurrentPage = 1;
            this.fetchRefundPaymentDetails();
        },

        // 重置退费明细筛选条件
        resetRefundPaymentDetailFilters() {
            this.refundPaymentDetailFilters = {
                id: '',
                student_id: '',
                order_id: '',
                refund_order_id: '',
                payment_id: '',
                payment_type: '',
                status: ''
            };
            this.fetchRefundPaymentDetails();
        },

        // 更改退费明细页码
        changeRefundPaymentDetailPage(page) {
            if (page >= 1 && page <= this.refundPaymentDetailTotalPages) {
                this.refundPaymentDetailCurrentPage = page;
            }
        },

        // 获取收款类型文本
        getPaymentTypeText(type) {
            const types = {
                0: '常规收款',
                1: '淘宝收款'
            };
            return types[type] || '未知';
        },

        // 获取收款主体文本
        getPayeeEntityText(entity) {
            if (entity === null || entity === undefined) {
                return '-';
            }
            const entities = {
                0: '北京',
                1: '西安'
            };
            return entities[entity] || '-';
        },

        // ==================== 审批流类型管理 ====================

        // 获取审批流类型列表
        async fetchApprovalFlowTypes() {
            this.loadingApprovalFlowTypes = true;
            try {
                const params = new URLSearchParams();
                if (this.approvalFlowTypeFilters.id) params.append('id', this.approvalFlowTypeFilters.id);
                if (this.approvalFlowTypeFilters.name) params.append('name', this.approvalFlowTypeFilters.name);
                if (this.approvalFlowTypeFilters.status !== '') params.append('status', this.approvalFlowTypeFilters.status);

                const response = await axios.get(`/api/approval-flow-types?${params.toString()}`, { withCredentials: true });
                this.approvalFlowTypes = response.data.data.approval_flow_types || [];
                this.filteredApprovalFlowTypes = this.approvalFlowTypes;
            } catch (err) {
                console.error('获取审批流类型失败:', err);
                alert(err.response?.data?.error || '获取审批流类型失败');
            } finally {
                this.loadingApprovalFlowTypes = false;
            }
        },

        // 搜索审批流类型
        searchApprovalFlowTypes() {
            this.fetchApprovalFlowTypes();
            this.approvalFlowTypeCurrentPage = 1;
        },

        // 重置审批流类型筛选
        resetApprovalFlowTypeFilters() {
            this.approvalFlowTypeFilters = {
                id: '',
                name: '',
                status: ''
            };
            this.fetchApprovalFlowTypes();
        },

        // 打开启用审批流类型确认
        openEnableFlowTypeConfirm(flowType) {
            this.currentFlowType = flowType;
            this.showEnableFlowTypeConfirm = true;
        },

        // 打开禁用审批流类型确认
        openDisableFlowTypeConfirm(flowType) {
            this.currentFlowType = flowType;
            this.showDisableFlowTypeConfirm = true;
        },

        // 启用审批流类型
        async doEnableFlowType() {
            try {
                await axios.put(`/api/approval-flow-types/${this.currentFlowType.id}/status`, {
                    status: 0
                }, { withCredentials: true });

                alert('审批流类型已启用');
                this.showEnableFlowTypeConfirm = false;
                this.currentFlowType = null;
                await this.fetchApprovalFlowTypes();
            } catch (err) {
                console.error('启用审批流类型失败:', err);
                alert(err.response?.data?.error || '启用审批流类型失败');
            }
        },

        // 禁用审批流类型
        async doDisableFlowType() {
            try {
                await axios.put(`/api/approval-flow-types/${this.currentFlowType.id}/status`, {
                    status: 1
                }, { withCredentials: true });

                alert('审批流类型已禁用');
                this.showDisableFlowTypeConfirm = false;
                this.currentFlowType = null;
                await this.fetchApprovalFlowTypes();
            } catch (err) {
                console.error('禁用审批流类型失败:', err);
                alert(err.response?.data?.error || '禁用审批流类型失败');
            }
        },

        // 审批流类型分页
        changeApprovalFlowTypePage(page) {
            if (page >= 1 && page <= this.approvalFlowTypeTotalPages) {
                this.approvalFlowTypeCurrentPage = page;
            }
        },

        // ==================== 审批流模板管理 ====================

        // 获取审批流模板列表
        async fetchApprovalFlowTemplates() {
            this.loadingApprovalFlowTemplates = true;
            try {
                const params = new URLSearchParams();
                if (this.approvalFlowTemplateFilters.id) params.append('id', this.approvalFlowTemplateFilters.id);
                if (this.approvalFlowTemplateFilters.approval_flow_type_id) params.append('approval_flow_type_id', this.approvalFlowTemplateFilters.approval_flow_type_id);
                if (this.approvalFlowTemplateFilters.name) params.append('name', this.approvalFlowTemplateFilters.name);
                if (this.approvalFlowTemplateFilters.status !== '') params.append('status', this.approvalFlowTemplateFilters.status);

                const response = await axios.get(`/api/approval-flow-templates?${params.toString()}`, { withCredentials: true });
                this.approvalFlowTemplates = response.data.data.approval_flow_templates || [];
                this.filteredApprovalFlowTemplates = this.approvalFlowTemplates;
            } catch (err) {
                console.error('获取审批流模板失败:', err);
                alert(err.response?.data?.error || '获取审批流模板失败');
            } finally {
                this.loadingApprovalFlowTemplates = false;
            }
        },

        // 搜索审批流模板
        searchApprovalFlowTemplates() {
            this.fetchApprovalFlowTemplates();
            this.approvalFlowTemplateCurrentPage = 1;
        },

        // 重置审批流模板筛选
        resetApprovalFlowTemplateFilters() {
            this.approvalFlowTemplateFilters = {
                id: '',
                approval_flow_type_id: '',
                name: '',
                status: ''
            };
            this.fetchApprovalFlowTemplates();
        },

        // 打开新增审批流模板弹窗
        async openAddApprovalFlowTemplateDrawer() {
            try {
                // 获取启用的审批流类型
                const typesResponse = await axios.get('/api/approval-flow-types?status=0', { withCredentials: true });
                this.activeApprovalFlowTypes = typesResponse.data.data.approval_flow_types || [];

                // 获取启用的用户账号
                const accountsResponse = await axios.get('/api/accounts?status=0', { withCredentials: true });
                this.activeUserAccounts = accountsResponse.data.accounts || [];

                this.showAddApprovalFlowTemplateDrawer = true;
            } catch (err) {
                console.error('获取数据失败:', err);
                alert(err.response?.data?.error || '获取数据失败');
            }
        },

        // 关闭新增审批流模板弹窗
        closeAddApprovalFlowTemplateDrawer() {
            this.showAddApprovalFlowTemplateDrawer = false;
            this.addApprovalFlowTemplateData = {
                name: '',
                approval_flow_type_id: '',
                nodes: [{
                    name: '',
                    type: 0,
                    approvers: ['']
                }],
                copy_users: []
            };
        },

        // 添加审批节点
        addApprovalNode() {
            this.addApprovalFlowTemplateData.nodes.push({
                name: '',
                type: 0,
                approvers: ['']
            });
        },

        // 删除审批节点
        removeApprovalNode(nodeIndex) {
            if (this.addApprovalFlowTemplateData.nodes.length > 1) {
                this.addApprovalFlowTemplateData.nodes.splice(nodeIndex, 1);
            }
        },

        // 添加审批人员
        addApprover(nodeIndex) {
            this.addApprovalFlowTemplateData.nodes[nodeIndex].approvers.push('');
        },

        // 删除审批人员
        removeApprover(nodeIndex, approverIndex) {
            const node = this.addApprovalFlowTemplateData.nodes[nodeIndex];
            if (node.approvers.length > 1) {
                node.approvers.splice(approverIndex, 1);
            }
        },

        // 添加抄送人员
        addCopyUser() {
            this.addApprovalFlowTemplateData.copy_users.push('');
        },

        // 删除抄送人员
        removeCopyUser(copyIndex) {
            this.addApprovalFlowTemplateData.copy_users.splice(copyIndex, 1);
        },

        // 保存审批流模板
        async saveApprovalFlowTemplate() {
            // 验证
            if (!this.addApprovalFlowTemplateData.name) {
                alert('请输入模板名称');
                return;
            }

            if (!this.addApprovalFlowTemplateData.approval_flow_type_id) {
                alert('请选择审批流类型');
                return;
            }

            // 验证节点
            for (let i = 0; i < this.addApprovalFlowTemplateData.nodes.length; i++) {
                const node = this.addApprovalFlowTemplateData.nodes[i];
                if (!node.name) {
                    alert(`请输入节点${i + 1}的名称`);
                    return;
                }

                // 过滤空的审批人员
                const validApprovers = node.approvers.filter(a => a !== '');
                if (validApprovers.length === 0) {
                    alert(`节点${i + 1}至少需要一个审批人员`);
                    return;
                }
                node.approvers = validApprovers;
            }

            // 过滤空的抄送人员
            this.addApprovalFlowTemplateData.copy_users = this.addApprovalFlowTemplateData.copy_users.filter(u => u !== '');

            try {
                const response = await axios.post('/api/approval-flow-templates', this.addApprovalFlowTemplateData, { withCredentials: true });

                alert('审批流模板创建成功');
                this.closeAddApprovalFlowTemplateDrawer();
                await this.fetchApprovalFlowTemplates();
            } catch (err) {
                console.error('创建审批流模板失败:', err);
                alert(err.response?.data?.error || '创建审批流模板失败');
            }
        },

        // 打开查看审批流模板详情弹窗
        async openViewApprovalFlowTemplateDrawer(template) {
            try {
                const response = await axios.get(`/api/approval-flow-templates/${template.id}`, { withCredentials: true });
                const data = response.data.data;

                this.viewApprovalFlowTemplateData = {
                    name: data.template.name,
                    approval_flow_type_id: data.template.approval_flow_type_id,
                    flow_type_name: data.template.flow_type_name,
                    nodes: data.nodes || [],
                    copy_users: data.copy_users || []
                };

                this.showViewApprovalFlowTemplateDrawer = true;
            } catch (err) {
                console.error('获取审批流模板详情失败:', err);
                alert(err.response?.data?.error || '获取审批流模板详情失败');
            }
        },

        // 关闭查看审批流模板详情弹窗
        closeViewApprovalFlowTemplateDrawer() {
            this.showViewApprovalFlowTemplateDrawer = false;
            this.viewApprovalFlowTemplateData = {
                name: '',
                approval_flow_type_id: '',
                flow_type_name: '',
                nodes: [],
                copy_users: []
            };
        },

        // 打开启用审批流模板确认
        openEnableFlowTemplateConfirm(template) {
            this.currentFlowTemplate = template;
            this.showEnableFlowTemplateConfirm = true;
        },

        // 打开禁用审批流模板确认
        openDisableFlowTemplateConfirm(template) {
            this.currentFlowTemplate = template;
            this.showDisableFlowTemplateConfirm = true;
        },

        // 启用审批流模板
        async doEnableFlowTemplate() {
            try {
                await axios.put(`/api/approval-flow-templates/${this.currentFlowTemplate.id}/status`, {
                    status: 0
                }, { withCredentials: true });

                alert('审批流模板已启用');
                this.showEnableFlowTemplateConfirm = false;
                this.currentFlowTemplate = null;
                await this.fetchApprovalFlowTemplates();
            } catch (err) {
                console.error('启用审批流模板失败:', err);
                alert(err.response?.data?.error || '启用审批流模板失败');
            }
        },

        // 禁用审批流模板
        async doDisableFlowTemplate() {
            try {
                await axios.put(`/api/approval-flow-templates/${this.currentFlowTemplate.id}/status`, {
                    status: 1
                }, { withCredentials: true });

                alert('审批流模板已禁用');
                this.showDisableFlowTemplateConfirm = false;
                this.currentFlowTemplate = null;
                await this.fetchApprovalFlowTemplates();
            } catch (err) {
                console.error('禁用审批流模板失败:', err);
                alert(err.response?.data?.error || '禁用审批流模板失败');
            }
        },

        // 审批流模板分页
        changeApprovalFlowTemplatePage(page) {
            if (page >= 1 && page <= this.approvalFlowTemplateTotalPages) {
                this.approvalFlowTemplateCurrentPage = page;
            }
        },

        // ==================== 审批流管理 ====================

        // 获取我发起的审批流列表
        async fetchInitiatedFlows() {
            this.loadingInitiatedFlows = true;
            try {
                const params = new URLSearchParams();
                if (this.initiatedFlowFilters.id) params.append('id', this.initiatedFlowFilters.id);
                if (this.initiatedFlowFilters.approval_flow_type_id) params.append('approval_flow_type_id', this.initiatedFlowFilters.approval_flow_type_id);
                if (this.initiatedFlowFilters.status !== '') params.append('status', this.initiatedFlowFilters.status);

                const response = await axios.get(`/api/approval-flows/initiated?${params.toString()}`, { withCredentials: true });
                this.initiatedFlows = response.data.data.initiated_flows || [];
            } catch (err) {
                console.error('获取我发起的审批流失败:', err);
                alert(err.response?.data?.error || '获取审批流列表失败');
            } finally {
                this.loadingInitiatedFlows = false;
            }
        },

        // 搜索我发起的审批流
        searchInitiatedFlows() {
            this.initiatedFlowCurrentPage = 1;
            this.fetchInitiatedFlows();
        },

        // 重置我发起的筛选条件
        resetInitiatedFlowFilters() {
            this.initiatedFlowFilters = { id: '', approval_flow_type_id: '', status: '' };
            this.fetchInitiatedFlows();
        },

        // 更改我发起的页码
        changeInitiatedFlowPage(page) {
            this.initiatedFlowCurrentPage = page;
        },

        // 打开撤销确认弹窗
        openCancelFlowConfirm(flow) {
            this.currentFlow = flow;
            this.showCancelFlowConfirm = true;
        },

        // 关闭撤销确认弹窗
        closeCancelFlowConfirm() {
            this.showCancelFlowConfirm = false;
            this.currentFlow = null;
        },

        // 撤销审批流
        async doCancelFlow() {
            try {
                await axios.put(`/api/approval-flows/${this.currentFlow.id}/cancel`, {}, { withCredentials: true });
                alert('撤销成功');
                this.closeCancelFlowConfirm();
                this.fetchInitiatedFlows();
            } catch (err) {
                console.error('撤销失败:', err);
                alert(err.response?.data?.error || '撤销失败');
            }
        },

        // 获取状态文本
        getFlowStatusText(status) {
            const statusMap = {
                0: '待审批',
                10: '已通过',
                20: '已驳回',
                99: '已撤销'
            };
            return statusMap[status] || '未知';
        },

        // 获取待我审批的列表
        async fetchPendingFlows() {
            this.loadingPendingFlows = true;
            try {
                const params = new URLSearchParams();
                if (this.pendingFlowFilters.id) params.append('id', this.pendingFlowFilters.id);
                if (this.pendingFlowFilters.approval_flow_id) params.append('approval_flow_id', this.pendingFlowFilters.approval_flow_id);
                if (this.pendingFlowFilters.approval_flow_type_id) params.append('approval_flow_type_id', this.pendingFlowFilters.approval_flow_type_id);

                const response = await axios.get(`/api/approval-flows/pending?${params.toString()}`, { withCredentials: true });
                this.pendingFlows = response.data.data.pending_flows || [];
            } catch (err) {
                console.error('获取待我审批列表失败:', err);
                alert(err.response?.data?.error || '获取列表失败');
            } finally {
                this.loadingPendingFlows = false;
            }
        },

        // 搜索待我审批
        searchPendingFlows() {
            this.pendingFlowCurrentPage = 1;
            this.fetchPendingFlows();
        },

        // 重置待我审批筛选条件
        resetPendingFlowFilters() {
            this.pendingFlowFilters = { id: '', approval_flow_id: '', approval_flow_type_id: '' };
            this.fetchPendingFlows();
        },

        // 更改待我审批页码
        changePendingFlowPage(page) {
            this.pendingFlowCurrentPage = page;
        },

        // 打开审批详情
        async openApprovalDetail(flow) {
            try {
                // 调用详情API获取完整信息
                const response = await axios.get(`/api/approval-flows/${flow.approval_flow_management_id}/detail`, { withCredentials: true });
                this.currentApprovalDetail = {
                    ...flow,
                    ...response.data.user_approval,
                    flow_info: response.data.flow_info,
                    all_nodes: response.data.all_nodes,
                    refund_order_info: response.data.refund_order_info
                };
                this.showApprovalDetailDrawer = true;
            } catch (err) {
                console.error('获取审批详情失败:', err);
                alert(err.response?.data?.error || '获取审批详情失败');
            }
        },

        // 关闭审批详情
        closeApprovalDetail() {
            this.showApprovalDetailDrawer = false;
            this.currentApprovalDetail = null;
        },

        // 审批通过
        async approveFlow() {
            try {
                await axios.post('/api/approval-flows/approve', {
                    node_case_user_id: this.currentApprovalDetail.id
                }, { withCredentials: true });
                alert('审批通过');
                this.closeApprovalDetail();
                this.fetchPendingFlows();
            } catch (err) {
                console.error('审批失败:', err);
                alert(err.response?.data?.error || '审批失败');
            }
        },

        // 审批驳回
        async rejectFlow() {
            try {
                await axios.post('/api/approval-flows/reject', {
                    node_case_user_id: this.currentApprovalDetail.id
                }, { withCredentials: true });
                alert('审批已驳回');
                this.closeApprovalDetail();
                this.fetchPendingFlows();
            } catch (err) {
                console.error('驳回失败:', err);
                alert(err.response?.data?.error || '驳回失败');
            }
        },

        // 获取处理完成的列表
        async fetchCompletedFlows() {
            this.loadingCompletedFlows = true;
            try {
                const params = new URLSearchParams();
                if (this.completedFlowFilters.id) params.append('id', this.completedFlowFilters.id);
                if (this.completedFlowFilters.approval_flow_id) params.append('approval_flow_id', this.completedFlowFilters.approval_flow_id);
                if (this.completedFlowFilters.approval_flow_type_id) params.append('approval_flow_type_id', this.completedFlowFilters.approval_flow_type_id);

                const response = await axios.get(`/api/approval-flows/completed?${params.toString()}`, { withCredentials: true });
                this.completedFlows = response.data.data.completed_flows || [];
            } catch (err) {
                console.error('获取处理完成列表失败:', err);
                alert(err.response?.data?.error || '获取列表失败');
            } finally {
                this.loadingCompletedFlows = false;
            }
        },

        // 搜索处理完成
        searchCompletedFlows() {
            this.completedFlowCurrentPage = 1;
            this.fetchCompletedFlows();
        },

        // 重置处理完成筛选条件
        resetCompletedFlowFilters() {
            this.completedFlowFilters = { id: '', approval_flow_id: '', approval_flow_type_id: '' };
            this.fetchCompletedFlows();
        },

        // 更改处理完成页码
        changeCompletedFlowPage(page) {
            this.completedFlowCurrentPage = page;
        },

        // 获取审批结果文本
        getApprovalResultText(result) {
            return result === 0 ? '通过' : '驳回';
        },

        // 获取抄送我的列表
        async fetchCopiedFlows() {
            this.loadingCopiedFlows = true;
            try {
                const params = new URLSearchParams();
                if (this.copiedFlowFilters.id) params.append('id', this.copiedFlowFilters.id);
                if (this.copiedFlowFilters.approval_flow_id) params.append('approval_flow_id', this.copiedFlowFilters.approval_flow_id);
                if (this.copiedFlowFilters.approval_flow_type_id) params.append('approval_flow_type_id', this.copiedFlowFilters.approval_flow_type_id);

                const response = await axios.get(`/api/approval-flows/copied?${params.toString()}`, { withCredentials: true });
                this.copiedFlows = response.data.data.copied_flows || [];
            } catch (err) {
                console.error('获取抄送列表失败:', err);
                alert(err.response?.data?.error || '获取列表失败');
            } finally {
                this.loadingCopiedFlows = false;
            }
        },

        // 搜索抄送我的
        searchCopiedFlows() {
            this.copiedFlowCurrentPage = 1;
            this.fetchCopiedFlows();
        },

        // 重置抄送我的筛选条件
        resetCopiedFlowFilters() {
            this.copiedFlowFilters = { id: '', approval_flow_id: '', approval_flow_type_id: '' };
            this.fetchCopiedFlows();
        },

        // 更改抄送我的页码
        changeCopiedFlowPage(page) {
            this.copiedFlowCurrentPage = page;
        },

        // ==================== 权限管理 ====================

        // 获取权限列表
        async fetchPermissions() {
            this.loadingPermissions = true;
            try {
                const params = new URLSearchParams();
                if (this.permissionFilters.id) params.append('id', this.permissionFilters.id);
                if (this.permissionFilters.menu_id) params.append('menu_id', this.permissionFilters.menu_id);

                const response = await axios.get(`/api/permissions?${params.toString()}`, { withCredentials: true });
                this.permissions = response.data.permissions || [];
                this.filteredPermissions = this.permissions;
            } catch (err) {
                console.error('获取权限列表失败:', err);
                alert(err.response?.data?.error || '获取权限列表失败');
            } finally {
                this.loadingPermissions = false;
            }
        },

        // 搜索权限
        searchPermissions() {
            this.fetchPermissions();
            this.permissionCurrentPage = 1;
        },

        // 重置权限筛选
        resetPermissionFilters() {
            this.permissionFilters = {
                id: '',
                menu_id: ''
            };
            this.fetchPermissions();
        },

        // 获取二级菜单列表
        async fetchSecondLevelMenus() {
            try {
                const response = await axios.get('/api/menu?level=2', { withCredentials: true });
                this.secondLevelMenus = response.data.menus || [];
            } catch (err) {
                console.error('获取菜单列表失败:', err);
            }
        },

        // 打开启用权限确认
        openEnablePermissionConfirm(permission) {
            this.currentPermission = permission;
            this.showEnablePermissionConfirm = true;
        },

        // 打开禁用权限确认
        openDisablePermissionConfirm(permission) {
            this.currentPermission = permission;
            this.showDisablePermissionConfirm = true;
        },

        // 启用权限
        async doEnablePermission() {
            try {
                await axios.put(`/api/permissions/${this.currentPermission.id}/status`, {
                    status: 0
                }, { withCredentials: true });

                alert('权限已启用');
                this.showEnablePermissionConfirm = false;
                this.currentPermission = null;
                await this.fetchPermissions();
                // 重新加载权限列表，使启用立即生效
                await this.fetchEnabledPermissions();
            } catch (err) {
                console.error('启用权限失败:', err);
                alert(err.response?.data?.error || '启用权限失败');
            }
        },

        // 禁用权限
        async doDisablePermission() {
            try {
                await axios.put(`/api/permissions/${this.currentPermission.id}/status`, {
                    status: 1
                }, { withCredentials: true });

                alert('权限已禁用');
                this.showDisablePermissionConfirm = false;
                this.currentPermission = null;
                await this.fetchPermissions();
                // 重新加载权限列表，使禁用立即生效
                await this.fetchEnabledPermissions();
            } catch (err) {
                console.error('禁用权限失败:', err);
                alert(err.response?.data?.error || '禁用权限失败');
            }
        },

        // 权限管理分页
        changePermissionPage(page) {
            if (page >= 1 && page <= this.permissionTotalPages) {
                this.permissionCurrentPage = page;
            }
        },

        // 同步当前用户的角色信息
        async syncRoleInfo() {
            try {
                const response = await axios.get('/api/sync-role', { withCredentials: true });

                if (response.data.data && response.data.data.role_changed) {
                    // 显示全局loading状态
                    this.syncingRole = true;

                    try {
                        // 角色发生变化，更新本地状态
                        this.isSuperAdmin = response.data.data.is_super_admin;
                        console.log('角色已更新:', response.data.data.is_super_admin ? '超级管理员' : '普通角色');

                        // 重新加载菜单和权限
                        await this.fetchMenus();
                        await this.fetchEnabledPermissions();

                        // 提示用户角色已更新
                        alert('您的角色权限已更新，页面已自动刷新');
                    } finally {
                        // 确保loading状态被关闭
                        this.syncingRole = false;
                    }
                }
            } catch (err) {
                // 同步失败时不影响正常使用，只记录错误
                console.error('同步角色信息失败:', err);
                this.syncingRole = false;
            }
        },

        // 获取当前用户的权限列表（用于按钮权限控制）
        async fetchEnabledPermissions() {
            this.loadingPermissions = true;
            try {
                const response = await axios.get('/api/user/permissions', { withCredentials: true });
                this.enabledPermissions = response.data.permissions || [];
                console.log('已加载用户权限列表:', this.enabledPermissions);
            } catch (err) {
                console.error('获取用户权限列表失败:', err);
                this.enabledPermissions = [];
            } finally {
                this.loadingPermissions = false;
            }
        },

        // 检查是否有某个权限
        hasPermission(actionId) {
            // 超级管理员拥有所有权限
            if (this.isSuperAdmin) {
                return true;
            }
            return this.enabledPermissions.includes(actionId);
        },

        // ==================== 属性管理功能 ====================

        // 获取属性数据
        async fetchAttributes() {
            this.loadingAttributes = true;
            try {
                const response = await axios.get('/api/attributes', { withCredentials: true });
                this.attributes = response.data.attributes;
                this.filteredAttributes = this.attributes;
            } catch (err) {
                console.error('获取属性失败:', err);
                this.error = '获取属性失败';
            } finally {
                this.loadingAttributes = false;
            }
        },

        // 属性筛选功能
        searchAttributes() {
            this.filteredAttributes = this.attributes.filter(attr => {
                const idMatch = !this.attributeFilters.id || attr.id === parseInt(this.attributeFilters.id);
                const classifyMatch = this.attributeFilters.classify === '' || attr.classify === parseInt(this.attributeFilters.classify);
                const statusMatch = this.attributeFilters.status === '' || attr.status === parseInt(this.attributeFilters.status);
                return idMatch && classifyMatch && statusMatch;
            });
        },

        // 重置属性筛选
        resetAttributeFilters() {
            this.attributeFilters = {
                id: '',
                classify: '',
                status: ''
            };
            this.filteredAttributes = this.attributes;
        },

        // 获取分类文本
        getClassifyText(classify) {
            return classify === 0 ? '属性' : '规格';
        },

        // 获取状态文本
        getAttributeStatusText(status) {
            return status === 0 ? '启用' : '禁用';
        },

        // 打开新增属性弹窗
        openAddAttributeModal() {
            this.showAddAttributeModal = true;
        },

        // 关闭新增属性弹窗
        closeAddAttributeModal() {
            this.showAddAttributeModal = false;
            this.addAttributeData = {
                name: '',
                classify: ''
            };
        },

        // 保存新增属性
        async saveAddAttribute() {
            if (!this.addAttributeData.name || this.addAttributeData.classify === '') {
                alert('请填写所有必填字段');
                return;
            }

            try {
                const response = await axios.post('/api/attributes', {
                    name: this.addAttributeData.name,
                    classify: parseInt(this.addAttributeData.classify)
                }, { withCredentials: true });

                if (response.data.message === '属性创建成功') {
                    await this.fetchAttributes();
                    this.closeAddAttributeModal();
                    alert('属性创建成功');
                }
            } catch (err) {
                console.error('创建属性失败:', err);
                alert(err.response?.data?.error || '创建属性失败');
            }
        },

        // 打开编辑属性弹窗
        openEditAttributeModal(attribute) {
            this.editAttributeData = {
                id: attribute.id,
                name: attribute.name,
                classify: attribute.classify
            };
            this.showEditAttributeModal = true;
        },

        // 关闭编辑属性弹窗
        closeEditAttributeModal() {
            this.showEditAttributeModal = false;
            this.editAttributeData = {
                id: '',
                name: '',
                classify: ''
            };
        },

        // 保存编辑属性
        async saveEditAttribute() {
            if (!this.editAttributeData.name || this.editAttributeData.classify === '') {
                alert('请填写所有必填字段');
                return;
            }

            try {
                const response = await axios.put(`/api/attributes/${this.editAttributeData.id}`, {
                    name: this.editAttributeData.name,
                    classify: parseInt(this.editAttributeData.classify)
                }, { withCredentials: true });

                if (response.data.message === '属性更新成功') {
                    await this.fetchAttributes();
                    this.closeEditAttributeModal();
                    alert('属性更新成功');
                }
            } catch (err) {
                console.error('更新属性失败:', err);
                alert(err.response?.data?.error || '更新属性失败');
            }
        },

        // 打开启用属性确认弹窗
        openEnableAttributeConfirm(attributeId) {
            this.attributeId = attributeId;
            this.showEnableAttributeConfirm = true;
        },

        // 关闭启用属性确认弹窗
        closeEnableAttributeConfirm() {
            this.showEnableAttributeConfirm = false;
            this.attributeId = null;
        },

        // 确认启用属性
        async confirmEnableAttribute() {
            try {
                const response = await axios.put(`/api/attributes/${this.attributeId}/status`, {
                    status: 0
                }, { withCredentials: true });

                if (response.data.message === '状态更新成功') {
                    await this.fetchAttributes();
                    this.closeEnableAttributeConfirm();
                    alert('属性已启用');
                }
            } catch (err) {
                console.error('启用属性失败:', err);
                alert(err.response?.data?.error || '启用属性失败');
            }
        },

        // 打开禁用属性确认弹窗
        openDisableAttributeConfirm(attributeId) {
            this.attributeId = attributeId;
            this.showDisableAttributeConfirm = true;
        },

        // 关闭禁用属性确认弹窗
        closeDisableAttributeConfirm() {
            this.showDisableAttributeConfirm = false;
            this.attributeId = null;
        },

        // 确认禁用属性
        async confirmDisableAttribute() {
            try {
                const response = await axios.put(`/api/attributes/${this.attributeId}/status`, {
                    status: 1
                }, { withCredentials: true });

                if (response.data.message === '状态更新成功') {
                    await this.fetchAttributes();
                    this.closeDisableAttributeConfirm();
                    alert('属性已禁用');
                }
            } catch (err) {
                console.error('禁用属性失败:', err);
                alert(err.response?.data?.error || '禁用属性失败');
            }
        },

        // ==================== 属性值管理功能 ====================

        // 打开属性值管理弹窗
        async openAttributeValuesModal(attribute) {
            this.currentAttributeId = attribute.id;
            this.currentAttributeName = attribute.name;

            // 获取该属性的属性值列表
            try {
                const response = await axios.get(`/api/attributes/${attribute.id}/values`, { withCredentials: true });
                const values = response.data.values;

                // 如果有属性值，则回显；否则显示一个空输入框
                if (values && values.length > 0) {
                    this.attributeValues = values.map(v => v.name);
                } else {
                    this.attributeValues = [''];
                }
            } catch (err) {
                console.error('获取属性值失败:', err);
                this.attributeValues = [''];
            }

            this.showAttributeValuesModal = true;
        },

        // 关闭属性值管理弹窗
        closeAttributeValuesModal() {
            this.showAttributeValuesModal = false;
            this.currentAttributeId = null;
            this.currentAttributeName = '';
            this.attributeValues = [''];
        },

        // 添加属性值输入框
        addAttributeValueInput() {
            this.attributeValues.push('');
        },

        // 删除属性值输入框
        removeAttributeValueInput(index) {
            if (this.attributeValues.length > 1) {
                this.attributeValues.splice(index, 1);
            }
        },

        // 清除属性值输入框内容（当仅剩一个时）
        clearAttributeValueInput(index) {
            this.attributeValues[index] = '';
        },

        // 保存属性值
        async saveAttributeValues() {
            // 过滤掉空的属性值
            const validValues = this.attributeValues.filter(v => v && v.trim());

            // 验证至少有一个属性值
            if (validValues.length === 0) {
                alert('至少需要填入一条属性值');
                return;
            }

            try {
                const response = await axios.post(`/api/attributes/${this.currentAttributeId}/values`, {
                    values: validValues
                }, { withCredentials: true });

                if (response.data.message === '属性值保存成功') {
                    alert('属性值保存成功');
                    await this.fetchAttributes();
                    this.closeAttributeValuesModal();
                }
            } catch (err) {
                console.error('保存属性值失败:', err);
                alert(err.response?.data?.error || '保存属性值失败');
            }
        },

        // ==================== 类型管理功能 ====================

        // 获取类型数据
        async fetchClassifies() {
            this.loadingClassifies = true;
            try {
                const response = await axios.get('/api/classifies', { withCredentials: true });
                this.classifies = response.data.classifies;
                this.filteredClassifies = this.classifies;
            } catch (err) {
                console.error('获取类型失败:', err);
                this.error = '获取类型失败';
            } finally {
                this.loadingClassifies = false;
            }
        },

        // 获取启用的一级类型列表（父级类型）
        async fetchParentClassifies() {
            try {
                const response = await axios.get('/api/classifies/parents', { withCredentials: true });
                this.parentClassifies = response.data.parents;
            } catch (err) {
                console.error('获取父级类型失败:', err);
                this.error = '获取父级类型失败';
            }
        },

        // 类型筛选功能
        searchClassifies() {
            this.filteredClassifies = this.classifies.filter(classify => {
                const idMatch = !this.classifyFilters.id || classify.id === parseInt(this.classifyFilters.id);
                const levelMatch = this.classifyFilters.level === '' || classify.level === parseInt(this.classifyFilters.level);
                const statusMatch = this.classifyFilters.status === '' || classify.status === parseInt(this.classifyFilters.status);
                return idMatch && levelMatch && statusMatch;
            });
        },

        // 重置类型筛选
        resetClassifyFilters() {
            this.classifyFilters = {
                id: '',
                level: '',
                status: ''
            };
            this.filteredClassifies = this.classifies;
        },

        // 获取级别文本
        getClassifyLevelText(level) {
            return level === 0 ? '一级' : '二级';
        },

        // 获取类型状态文本
        getClassifyStatusText(status) {
            return status === 0 ? '启用' : '禁用';
        },

        // 打开新增类型弹窗
        async openAddClassifyModal() {
            await this.fetchParentClassifies();
            this.showAddClassifyModal = true;
        },

        // 关闭新增类型弹窗
        closeAddClassifyModal() {
            this.showAddClassifyModal = false;
            this.addClassifyData = {
                name: '',
                level: '',
                parent_id: ''
            };
        },

        // 新增类型级别变化时的处理
        onAddClassifyLevelChange() {
            // 如果切换到一级，清空父级选择
            if (this.addClassifyData.level === '0') {
                this.addClassifyData.parent_id = '';
            }
        },

        // 保存新增类型
        async saveAddClassify() {
            if (!this.addClassifyData.name || this.addClassifyData.level === '') {
                alert('请填写所有必填字段');
                return;
            }

            // 如果是二级类型，必须选择父级
            if (this.addClassifyData.level === '1' && !this.addClassifyData.parent_id) {
                alert('请选择父级类型');
                return;
            }

            try {
                const response = await axios.post('/api/classifies', {
                    name: this.addClassifyData.name,
                    level: parseInt(this.addClassifyData.level),
                    parent_id: this.addClassifyData.level === '1' ? parseInt(this.addClassifyData.parent_id) : null
                }, { withCredentials: true });

                if (response.data.message === '类型创建成功') {
                    await this.fetchClassifies();
                    this.closeAddClassifyModal();
                    alert('类型创建成功');
                }
            } catch (err) {
                console.error('创建类型失败:', err);
                alert(err.response?.data?.error || '创建类型失败');
            }
        },

        // 打开编辑类型弹窗
        async openEditClassifyModal(classify) {
            await this.fetchParentClassifies();
            this.editClassifyData = {
                id: classify.id,
                name: classify.name,
                level: classify.level,
                parent_id: classify.parent_id || ''
            };
            this.showEditClassifyModal = true;
        },

        // 关闭编辑类型弹窗
        closeEditClassifyModal() {
            this.showEditClassifyModal = false;
            this.editClassifyData = {
                id: '',
                name: '',
                level: '',
                parent_id: ''
            };
        },

        // 编辑类型级别变化时的处理
        onEditClassifyLevelChange() {
            // 如果切换到一级，清空父级选择
            const level = typeof this.editClassifyData.level === 'string'
                ? parseInt(this.editClassifyData.level)
                : this.editClassifyData.level;
            if (level === 0) {
                this.editClassifyData.parent_id = '';
            }
        },

        // 保存编辑类型
        async saveEditClassify() {
            const level = typeof this.editClassifyData.level === 'string'
                ? parseInt(this.editClassifyData.level)
                : this.editClassifyData.level;

            if (!this.editClassifyData.name || this.editClassifyData.level === '') {
                alert('请填写所有必填字段');
                return;
            }

            // 如果是二级类型，必须选择父级
            if (level === 1 && !this.editClassifyData.parent_id) {
                alert('请选择父级类型');
                return;
            }

            try {
                const response = await axios.put(`/api/classifies/${this.editClassifyData.id}`, {
                    name: this.editClassifyData.name,
                    level: level,
                    parent_id: level === 1 ? parseInt(this.editClassifyData.parent_id) : null
                }, { withCredentials: true });

                if (response.data.message === '类型更新成功') {
                    await this.fetchClassifies();
                    this.closeEditClassifyModal();
                    alert('类型更新成功');
                }
            } catch (err) {
                console.error('更新类型失败:', err);
                alert(err.response?.data?.error || '更新类型失败');
            }
        },

        // 打开启用类型确认弹窗
        openEnableClassifyConfirm(classifyId) {
            this.classifyId = classifyId;
            this.showEnableClassifyConfirm = true;
        },

        // 关闭启用类型确认弹窗
        closeEnableClassifyConfirm() {
            this.showEnableClassifyConfirm = false;
            this.classifyId = null;
        },

        // 确认启用类型
        async confirmEnableClassify() {
            try {
                const response = await axios.put(`/api/classifies/${this.classifyId}/status`, {
                    status: 0
                }, { withCredentials: true });

                if (response.data.message === '状态更新成功') {
                    await this.fetchClassifies();
                    this.closeEnableClassifyConfirm();
                    alert('类型已启用');
                }
            } catch (err) {
                console.error('启用类型失败:', err);
                alert(err.response?.data?.error || '启用类型失败');
            }
        },

        // 打开禁用类型确认弹窗
        openDisableClassifyConfirm(classifyId) {
            this.classifyId = classifyId;
            this.showDisableClassifyConfirm = true;
        },

        // 关闭禁用类型确认弹窗
        closeDisableClassifyConfirm() {
            this.showDisableClassifyConfirm = false;
            this.classifyId = null;
        },

        // 确认禁用类型
        async confirmDisableClassify() {
            try {
                const response = await axios.put(`/api/classifies/${this.classifyId}/status`, {
                    status: 1
                }, { withCredentials: true });

                if (response.data.message === '状态更新成功') {
                    await this.fetchClassifies();
                    this.closeDisableClassifyConfirm();
                    alert('类型已禁用');
                }
            } catch (err) {
                console.error('禁用类型失败:', err);
                alert(err.response?.data?.error || '禁用类型失败');
            }
        },

        // ==================== 分页功能 ====================

        // 学生分页
        changeStudentPage(page) {
            if (page >= 1 && page <= this.studentTotalPages) {
                this.studentCurrentPage = page;
            }
        },
        // 教练分页
        changeCoachPage(page) {
            if (page >= 1 && page <= this.coachTotalPages) {
                this.coachCurrentPage = page;
            }
        },
        // 订单分页
        changeOrderPage(page) {
            if (page >= 1 && page <= this.orderTotalPages) {
                this.orderCurrentPage = page;
            }
        },

        // ==================== 子订单管理方法 ====================

        // 获取子订单数据
        async fetchChildOrders() {
            this.loadingChildOrders = true;
            try {
                const response = await axios.get('/api/childorders', { withCredentials: true });
                this.childOrders = response.data.childorders;
                this.filteredChildOrders = this.childOrders;
            } catch (err) {
                console.error('获取子订单数据失败:', err);
                this.error = '获取子订单数据失败';
            } finally {
                this.loadingChildOrders = false;
            }
        },

        // 搜索子订单
        searchChildOrders() {
            this.filteredChildOrders = this.childOrders.filter(order => {
                let match = true;
                if (this.childOrderFilters.id && order.id !== parseInt(this.childOrderFilters.id)) {
                    match = false;
                }
                if (this.childOrderFilters.parentsid && order.parentsid !== parseInt(this.childOrderFilters.parentsid)) {
                    match = false;
                }
                if (this.childOrderFilters.goodsid && order.goodsid !== parseInt(this.childOrderFilters.goodsid)) {
                    match = false;
                }
                if (this.childOrderFilters.status && order.status !== parseInt(this.childOrderFilters.status)) {
                    match = false;
                }
                return match;
            });
            this.childOrderCurrentPage = 1;
        },

        // 重置子订单筛选
        resetChildOrderFilters() {
            this.childOrderFilters = {
                id: '',
                parentsid: '',
                goodsid: '',
                status: ''
            };
            this.filteredChildOrders = this.childOrders;
            this.childOrderCurrentPage = 1;
        },

        // 获取子订单状态文本
        getChildOrderStatusText(status) {
            const statusMap = {
                0: '草稿',
                10: '未支付',
                20: '部分支付',
                30: '已支付',
                99: '已作废'
            };
            return statusMap[status] || '未知';
        },

        // 子订单分页
        changeChildOrderPage(page) {
            if (page >= 1 && page <= this.childOrderTotalPages) {
                this.childOrderCurrentPage = page;
            }
        },

        // 账号分页
        changeAccountPage(page) {
            if (page >= 1 && page <= this.accountTotalPages) {
                this.accountCurrentPage = page;
            }
        },
        // 属性分页
        changeAttributePage(page) {
            if (page >= 1 && page <= this.attributeTotalPages) {
                this.attributeCurrentPage = page;
            }
        },
        // 类型分页
        changeClassifyPage(page) {
            if (page >= 1 && page <= this.classifyTotalPages) {
                this.classifyCurrentPage = page;
            }
        },

        // ==================== 品牌管理功能 ====================

        // 获取品牌数据
        async fetchBrands() {
            this.loadingBrands = true;
            try {
                const response = await axios.get('/api/brands', { withCredentials: true });
                this.brands = response.data.brands;
                this.filteredBrands = this.brands;
            } catch (err) {
                console.error('获取品牌失败:', err);
                this.error = '获取品牌失败';
            } finally {
                this.loadingBrands = false;
            }
        },

        // 品牌筛选功能
        searchBrands() {
            this.filteredBrands = this.brands.filter(brand => {
                const idMatch = !this.brandFilters.id || brand.id === parseInt(this.brandFilters.id);
                const nameMatch = !this.brandFilters.name || brand.name === this.brandFilters.name;
                const statusMatch = this.brandFilters.status === '' || brand.status === parseInt(this.brandFilters.status);
                return idMatch && nameMatch && statusMatch;
            });
            // 筛选后重置到第一页
            this.brandCurrentPage = 1;
        },

        // 重置品牌筛选
        resetBrandFilters() {
            this.brandFilters = {
                id: '',
                name: '',
                status: ''
            };
            this.filteredBrands = this.brands;
            this.brandCurrentPage = 1;
        },

        // 获取品牌状态文本
        getBrandStatusText(status) {
            return status === 0 ? '启用' : '禁用';
        },

        // 打开新增品牌弹窗
        openAddBrandModal() {
            this.showAddBrandModal = true;
        },

        // 关闭新增品牌弹窗
        closeAddBrandModal() {
            this.showAddBrandModal = false;
            this.addBrandData = {
                name: ''
            };
        },

        // 保存新增品牌
        async saveAddBrand() {
            if (!this.addBrandData.name) {
                alert('请填写品牌名称');
                return;
            }

            try {
                const response = await axios.post('/api/brands', {
                    name: this.addBrandData.name
                }, { withCredentials: true });

                if (response.data.message === '品牌添加成功') {
                    await this.fetchBrands();
                    this.closeAddBrandModal();
                    alert('品牌添加成功');
                }
            } catch (err) {
                console.error('新增品牌失败:', err);
                alert(err.response?.data?.error || '新增品牌失败');
            }
        },

        // 打开编辑品牌弹窗
        openEditBrandModal(brand) {
            this.editBrandData = {
                id: brand.id,
                name: brand.name
            };
            this.showEditBrandModal = true;
        },

        // 关闭编辑品牌弹窗
        closeEditBrandModal() {
            this.showEditBrandModal = false;
            this.editBrandData = {
                id: '',
                name: ''
            };
        },

        // 保存编辑品牌
        async saveEditBrand() {
            if (!this.editBrandData.name) {
                alert('请填写品牌名称');
                return;
            }

            try {
                const response = await axios.put(`/api/brands/${this.editBrandData.id}`, {
                    name: this.editBrandData.name
                }, { withCredentials: true });

                if (response.data.message === '品牌信息更新成功') {
                    await this.fetchBrands();
                    this.closeEditBrandModal();
                    alert('品牌信息更新成功');
                }
            } catch (err) {
                console.error('更新品牌失败:', err);
                alert(err.response?.data?.error || '更新品牌失败');
            }
        },

        // 打开启用品牌确认弹窗
        openEnableBrandConfirm(brandId) {
            this.brandId = brandId;
            this.showEnableBrandConfirm = true;
        },

        // 关闭启用品牌确认弹窗
        closeEnableBrandConfirm() {
            this.showEnableBrandConfirm = false;
            this.brandId = null;
        },

        // 确认启用品牌
        async confirmEnableBrand() {
            try {
                const response = await axios.put(`/api/brands/${this.brandId}/status`, {
                    status: 0
                }, { withCredentials: true });

                if (response.data.message === '品牌状态更新成功') {
                    await this.fetchBrands();
                    this.closeEnableBrandConfirm();
                    alert('品牌已启用');
                }
            } catch (err) {
                console.error('启用品牌失败:', err);
                alert(err.response?.data?.error || '启用品牌失败');
            }
        },

        // 打开禁用品牌确认弹窗
        openDisableBrandConfirm(brandId) {
            this.brandId = brandId;
            this.showDisableBrandConfirm = true;
        },

        // 关闭禁用品牌确认弹窗
        closeDisableBrandConfirm() {
            this.showDisableBrandConfirm = false;
            this.brandId = null;
        },

        // 确认禁用品牌
        async confirmDisableBrand() {
            try {
                const response = await axios.put(`/api/brands/${this.brandId}/status`, {
                    status: 1
                }, { withCredentials: true });

                if (response.data.message === '品牌状态更新成功') {
                    await this.fetchBrands();
                    this.closeDisableBrandConfirm();
                    alert('品牌已禁用');
                }
            } catch (err) {
                console.error('禁用品牌失败:', err);
                alert(err.response?.data?.error || '禁用品牌失败');
            }
        },

        // 品牌分页
        changeBrandPage(page) {
            if (page >= 1 && page <= this.brandTotalPages) {
                this.brandCurrentPage = page;
            }
        },

        // ==================== 商品管理功能 ====================

        // 格式化商品属性显示（只显示第一个属性和第一个规格）
        formatGoodsAttributes(goods) {
            if (!goods.attributes_full || goods.attributes_full.length === 0) {
                return '';
            }

            // 如果只有1-2个，直接显示
            if (goods.attributes_full.length <= 2) {
                return goods.attributes_full.join(',');
            }

            // 否则只显示前2个
            return goods.attributes_full.slice(0, 2).join(',') + '...';
        },

        // 获取完整属性列表用于悬浮显示
        getFullAttributes(goods) {
            if (!goods.attributes_full || goods.attributes_full.length === 0) {
                return '';
            }
            return goods.attributes_full.join(',');
        },

        // 获取商品数据
        async fetchGoods() {
            this.loadingGoods = true;
            try {
                const response = await axios.get('/api/goods', { withCredentials: true });
                this.goods = response.data.goods;
                this.filteredGoods = this.goods;
            } catch (err) {
                console.error('获取商品失败:', err);
                this.error = '获取商品失败';
            } finally {
                this.loadingGoods = false;
            }
        },

        // 商品筛选功能
        searchGoods() {
            this.filteredGoods = this.goods.filter(goods => {
                const idMatch = !this.goodsFilters.id || goods.id === parseInt(this.goodsFilters.id);
                const nameMatch = !this.goodsFilters.name || goods.name.toLowerCase().includes(this.goodsFilters.name.toLowerCase());
                const brandMatch = this.goodsFilters.brandid === '' || goods.brandid === parseInt(this.goodsFilters.brandid);
                const classifyMatch = this.goodsFilters.classifyid === '' || goods.classifyid === parseInt(this.goodsFilters.classifyid);
                const statusMatch = this.goodsFilters.status === '' || goods.status === parseInt(this.goodsFilters.status);
                return idMatch && nameMatch && brandMatch && classifyMatch && statusMatch;
            });
            this.goodsCurrentPage = 1;
        },

        // 重置商品筛选
        resetGoodsFilters() {
            this.goodsFilters = {
                id: '',
                name: '',
                brandid: '',
                classifyid: '',
                status: ''
            };
            this.filteredGoods = this.goods;
            this.goodsCurrentPage = 1;
        },

        // 打开新增商品抽屉
        async openAddGoodsDrawer() {
            // 获取启用的品牌列表
            try {
                const response = await axios.get('/api/brands/active', { withCredentials: true });
                this.activeBrands = response.data.brands;
            } catch (err) {
                console.error('获取启用品牌列表失败:', err);
                this.error = '获取启用品牌列表失败';
            }

            // 获取启用的类型列表
            try {
                const response = await axios.get('/api/classifies/active', { withCredentials: true });
                this.activeClassifies = response.data.classifies;
            } catch (err) {
                console.error('获取启用类型列表失败:', err);
                this.error = '获取启用类型列表失败';
            }

            // 获取启用的属性列表（包含属性值）
            try {
                const response = await axios.get('/api/attributes/active', { withCredentials: true });
                this.activeAttributes = response.data.attributes;
            } catch (err) {
                console.error('获取启用属性列表失败:', err);
                this.error = '获取启用属性列表失败';
            }

            this.showAddGoodsDrawer = true;
        },

        // 关闭新增商品抽屉
        closeAddGoodsDrawer() {
            this.showAddGoodsDrawer = false;
            this.addGoodsData = {
                name: '',
                brandid: '',
                classifyid: '',
                isgroup: 1,
                price: '',
                attributevalue_ids: []
            };
            this.selectedIncludedGoods = [];
            this.goodsAttributeRows = [{ attributeId: '', valueId: '' }];
            this.goodsSpecRows = [{ attributeId: '', valueId: '' }];
        },

        // 获取某行可用的属性列表（排除已选择的属性，只返回classify=0的）
        getAvailableAttributes(currentIndex) {
            const selectedAttrIds = this.goodsAttributeRows
                .map((row, index) => index !== currentIndex ? String(row.attributeId) : null)
                .filter(id => id && id !== '');
            return this.attributesList.filter(attr => !selectedAttrIds.includes(String(attr.id)));
        },

        // 获取某行可用的规格列表（排除已选择的规格，只返回classify=1的）
        getAvailableSpecs(currentIndex) {
            const selectedSpecIds = this.goodsSpecRows
                .map((row, index) => index !== currentIndex ? String(row.attributeId) : null)
                .filter(id => id && id !== '');
            return this.specsList.filter(spec => !selectedSpecIds.includes(String(spec.id)));
        },

        // 获取某个属性的属性值列表
        getAttributeValues(attributeId) {
            if (!attributeId) return [];
            const attribute = this.activeAttributes.find(attr => attr.id === parseInt(attributeId));
            return attribute && attribute.values ? attribute.values : [];
        },

        // 当某行的属性选择变化时，清空该行的属性值
        onRowAttributeChange(index) {
            this.goodsAttributeRows[index].valueId = '';
        },

        // 添加新的属性行
        addAttributeRow() {
            this.goodsAttributeRows.push({ attributeId: '', valueId: '' });
        },

        // 删除属性行
        removeAttributeRow(index) {
            if (this.goodsAttributeRows.length > 1) {
                this.goodsAttributeRows.splice(index, 1);
            }
        },

        // 当某行的规格选择变化时，清空该行的规格值
        onRowSpecChange(index) {
            this.goodsSpecRows[index].valueId = '';
        },

        // 添加新的规格行
        addSpecRow() {
            this.goodsSpecRows.push({ attributeId: '', valueId: '' });
        },

        // 删除规格行
        removeSpecRow(index) {
            if (this.goodsSpecRows.length > 1) {
                this.goodsSpecRows.splice(index, 1);
            }
        },

        // 组合售卖变化处理
        async onIsGroupChange() {
            // 如果切换为组合售卖，加载可用商品列表
            if (this.addGoodsData.isgroup == 0 || this.editGoodsData.isgroup == 0) {
                try {
                    const excludeId = this.editGoodsData.id || null;
                    const response = await axios.get('/api/goods/available-for-combo', {
                        params: excludeId ? { exclude_id: excludeId } : {},
                        withCredentials: true
                    });
                    this.availableGoodsForCombo = response.data.goods;
                } catch (err) {
                    console.error('获取可用商品列表失败:', err);
                    alert('获取可用商品列表失败');
                }
            } else {
                // 切换为非组合时，清空包含商品列表
                this.selectedIncludedGoods = [];
            }
        },

        // 打开新增包含商品子弹窗
        openAddIncludedGoodsModal() {
            this.addIncludedGoodsData = {
                goods_id: '',
                goods_name: '',
                brand_name: '',
                classify_name: '',
                attributes: '',
                price: ''
            };
            this.showAddIncludedGoodsModal = true;
        },

        // 包含商品选择变化
        onIncludedGoodsChange() {
            const selectedGoods = this.availableGoodsForCombo.find(
                g => g.id === parseInt(this.addIncludedGoodsData.goods_id)
            );

            if (selectedGoods) {
                this.addIncludedGoodsData.goods_name = selectedGoods.name;
                this.addIncludedGoodsData.brand_name = selectedGoods.brand_name;
                this.addIncludedGoodsData.classify_name = selectedGoods.classify_name;
                this.addIncludedGoodsData.attributes = selectedGoods.attributes;
                this.addIncludedGoodsData.price = selectedGoods.price;
            }
        },

        // 保存包含商品
        saveIncludedGoods() {
            if (!this.addIncludedGoodsData.goods_id) {
                alert('请选择商品');
                return;
            }

            // 检查是否已添加
            const exists = this.selectedIncludedGoods.find(
                g => g.id === parseInt(this.addIncludedGoodsData.goods_id)
            );
            if (exists) {
                alert('该商品已添加，请勿重复添加');
                return;
            }

            // 添加到列表
            const selectedGoods = this.availableGoodsForCombo.find(
                g => g.id === parseInt(this.addIncludedGoodsData.goods_id)
            );
            if (selectedGoods) {
                this.selectedIncludedGoods.push({
                    id: selectedGoods.id,
                    name: selectedGoods.name,
                    brand_name: selectedGoods.brand_name,
                    classify_name: selectedGoods.classify_name,
                    attributes: selectedGoods.attributes,
                    price: selectedGoods.price
                });
            }

            this.closeAddIncludedGoodsModal();
        },

        // 关闭包含商品子弹窗
        closeAddIncludedGoodsModal() {
            this.showAddIncludedGoodsModal = false;
            this.addIncludedGoodsData = {
                goods_id: '',
                goods_name: '',
                brand_name: '',
                classify_name: '',
                attributes: '',
                price: ''
            };
        },

        // 删除包含商品
        removeIncludedGoods(index) {
            this.selectedIncludedGoods.splice(index, 1);
        },

        // 保存新增商品
        async saveAddGoods() {
            if (!this.addGoodsData.name || !this.addGoodsData.brandid ||
                !this.addGoodsData.classifyid || !this.addGoodsData.price) {
                alert('请填写所有必填字段');
                return;
            }

            // 组合商品验证
            if (this.addGoodsData.isgroup == 0 && this.selectedIncludedGoods.length === 0) {
                alert('组合商品必须至少包含一个子商品');
                return;
            }

            // 从goodsAttributeRows和goodsSpecRows中提取属性值ID（过滤掉空值）
            const attributeValueIds = this.goodsAttributeRows
                .filter(row => row.attributeId && row.valueId)
                .map(row => parseInt(row.valueId));

            const specValueIds = this.goodsSpecRows
                .filter(row => row.attributeId && row.valueId)
                .map(row => parseInt(row.valueId));

            // 合并属性值和规格值
            const attributevalue_ids = [...attributeValueIds, ...specValueIds];

            // 提取包含商品ID
            const included_goods_ids = this.selectedIncludedGoods.map(g => g.id);

            console.log('goodsAttributeRows:', this.goodsAttributeRows);
            console.log('goodsSpecRows:', this.goodsSpecRows);
            console.log('提取的属性值IDs:', attributevalue_ids);
            console.log('包含商品IDs:', included_goods_ids);

            try {
                const response = await axios.post('/api/goods', {
                    name: this.addGoodsData.name,
                    brandid: parseInt(this.addGoodsData.brandid),
                    classifyid: parseInt(this.addGoodsData.classifyid),
                    isgroup: parseInt(this.addGoodsData.isgroup),
                    price: parseFloat(this.addGoodsData.price),
                    attributevalue_ids: attributevalue_ids,
                    included_goods_ids: included_goods_ids
                }, { withCredentials: true });

                if (response.data.message === '商品添加成功') {
                    await this.fetchGoods();
                    this.closeAddGoodsDrawer();
                    alert('商品添加成功');
                }
            } catch (err) {
                console.error('新增商品失败:', err);
                alert(err.response?.data?.error || '新增商品失败');
            }
        },

        // 打开编辑商品抽屉
        async openEditGoodsDrawer(goodsId) {
            // 获取启用的品牌列表
            try {
                const response = await axios.get('/api/brands/active', { withCredentials: true });
                this.activeBrands = response.data.brands;
            } catch (err) {
                console.error('获取启用品牌列表失败:', err);
                this.error = '获取启用品牌列表失败';
            }

            // 获取启用的类型列表
            try {
                const response = await axios.get('/api/classifies/active', { withCredentials: true });
                this.activeClassifies = response.data.classifies;
            } catch (err) {
                console.error('获取启用类型列表失败:', err);
                this.error = '获取启用类型列表失败';
            }

            // 获取启用的属性列表（包含属性值）
            try {
                const response = await axios.get('/api/attributes/active', { withCredentials: true });
                this.activeAttributes = response.data.attributes;
            } catch (err) {
                console.error('获取启用属性列表失败:', err);
                this.error = '获取启用属性列表失败';
            }

            // 获取商品详情
            try {
                const response = await axios.get(`/api/goods/${goodsId}`, { withCredentials: true });
                const goods = response.data.goods;
                this.editGoodsData = {
                    id: goods.id,
                    name: goods.name,
                    brandid: goods.brandid,
                    classifyid: goods.classifyid,
                    isgroup: goods.isgroup,
                    price: goods.price,
                    attributevalue_ids: goods.attributevalue_ids || []
                };

                // 根据商品的属性值ID重建goodsAttributeRows和goodsSpecRows数组
                this.goodsAttributeRows = [];
                this.goodsSpecRows = [];
                if (goods.attributevalue_ids && goods.attributevalue_ids.length > 0) {
                    for (const valueId of goods.attributevalue_ids) {
                        // 找到该属性值所属的属性
                        for (const attribute of this.activeAttributes) {
                            const value = attribute.values.find(v => v.id === valueId);
                            if (value) {
                                const row = {
                                    attributeId: attribute.id.toString(),
                                    valueId: value.id.toString()
                                };
                                // 根据classify区分是属性还是规格
                                if (attribute.classify === 0) {
                                    this.goodsAttributeRows.push(row);
                                } else if (attribute.classify === 1) {
                                    this.goodsSpecRows.push(row);
                                }
                                break;
                            }
                        }
                    }
                }
                // 如果没有属性值，至少保留一个空行
                if (this.goodsAttributeRows.length === 0) {
                    this.goodsAttributeRows = [{ attributeId: '', valueId: '' }];
                }
                if (this.goodsSpecRows.length === 0) {
                    this.goodsSpecRows = [{ attributeId: '', valueId: '' }];
                }

                // 如果是组合商品，加载包含商品数据
                if (goods.isgroup == 0) {
                    // 加载可用商品列表（用于子弹窗）
                    try {
                        const availableResponse = await axios.get('/api/goods/available-for-combo', {
                            params: { exclude_id: goodsId },
                            withCredentials: true
                        });
                        this.availableGoodsForCombo = availableResponse.data.goods;
                    } catch (err) {
                        console.error('获取可用商品列表失败:', err);
                    }

                    // 加载已包含的商品列表
                    try {
                        const includedResponse = await axios.get(`/api/goods/${goodsId}/included-goods`, {
                            withCredentials: true
                        });
                        this.selectedIncludedGoods = includedResponse.data.included_goods;
                    } catch (err) {
                        console.error('获取包含商品列表失败:', err);
                    }
                } else {
                    this.selectedIncludedGoods = [];
                }
            } catch (err) {
                console.error('获取商品详情失败:', err);
                alert('获取商品详情失败');
                return;
            }

            this.showEditGoodsDrawer = true;
        },

        // 关闭编辑商品抽屉
        closeEditGoodsDrawer() {
            this.showEditGoodsDrawer = false;
            this.editGoodsData = {
                id: '',
                name: '',
                brandid: '',
                classifyid: '',
                isgroup: 1,
                price: '',
                attributevalue_ids: []
            };
            this.selectedIncludedGoods = [];
            this.goodsAttributeRows = [{ attributeId: '', valueId: '' }];
            this.goodsSpecRows = [{ attributeId: '', valueId: '' }];
        },

        // 保存编辑商品
        async saveEditGoods() {
            if (!this.editGoodsData.name || !this.editGoodsData.brandid ||
                !this.editGoodsData.classifyid || !this.editGoodsData.price) {
                alert('请填写所有必填字段');
                return;
            }

            // 组合商品验证
            if (this.editGoodsData.isgroup == 0 && this.selectedIncludedGoods.length === 0) {
                alert('组合商品必须至少包含一个子商品');
                return;
            }

            // 从goodsAttributeRows和goodsSpecRows中提取属性值ID（过滤掉空值）
            const attributeValueIds = this.goodsAttributeRows
                .filter(row => row.attributeId && row.valueId)
                .map(row => parseInt(row.valueId));

            const specValueIds = this.goodsSpecRows
                .filter(row => row.attributeId && row.valueId)
                .map(row => parseInt(row.valueId));

            // 合并属性值和规格值
            const attributevalue_ids = [...attributeValueIds, ...specValueIds];

            // 提取包含商品ID
            const included_goods_ids = this.selectedIncludedGoods.map(g => g.id);

            try {
                const response = await axios.put(`/api/goods/${this.editGoodsData.id}`, {
                    name: this.editGoodsData.name,
                    brandid: parseInt(this.editGoodsData.brandid),
                    classifyid: parseInt(this.editGoodsData.classifyid),
                    price: parseFloat(this.editGoodsData.price),
                    attributevalue_ids: attributevalue_ids,
                    included_goods_ids: included_goods_ids
                }, { withCredentials: true });

                if (response.data.message === '商品信息更新成功') {
                    await this.fetchGoods();
                    this.closeEditGoodsDrawer();
                    alert('商品信息更新成功');
                }
            } catch (err) {
                console.error('更新商品失败:', err);
                alert(err.response?.data?.error || '更新商品失败');
            }
        },

        // 打开启用商品确认弹窗
        openEnableGoodsConfirm(goodsId) {
            this.goodsId = goodsId;
            this.showEnableGoodsConfirm = true;
        },

        // 关闭启用商品确认弹窗
        closeEnableGoodsConfirm() {
            this.showEnableGoodsConfirm = false;
            this.goodsId = null;
        },

        // 确认启用商品
        async confirmEnableGoods() {
            try {
                const response = await axios.put(`/api/goods/${this.goodsId}/status`, {
                    status: 0
                }, { withCredentials: true });

                if (response.data.message === '商品状态更新成功') {
                    await this.fetchGoods();
                    this.closeEnableGoodsConfirm();
                    alert('商品已启用');
                }
            } catch (err) {
                console.error('启用商品失败:', err);
                alert(err.response?.data?.error || '启用商品失败');
            }
        },

        // 打开禁用商品确认弹窗
        openDisableGoodsConfirm(goodsId) {
            this.goodsId = goodsId;
            this.showDisableGoodsConfirm = true;
        },

        // 关闭禁用商品确认弹窗
        closeDisableGoodsConfirm() {
            this.showDisableGoodsConfirm = false;
            this.goodsId = null;
        },

        // 确认禁用商品
        async confirmDisableGoods() {
            try {
                const response = await axios.put(`/api/goods/${this.goodsId}/status`, {
                    status: 1
                }, { withCredentials: true });

                if (response.data.message === '商品状态更新成功') {
                    await this.fetchGoods();
                    this.closeDisableGoodsConfirm();
                    alert('商品已禁用');
                }
            } catch (err) {
                console.error('禁用商品失败:', err);
                alert(err.response?.data?.error || '禁用商品失败');
            }
        },

        // 商品分页
        changeGoodsPage(page) {
            if (page >= 1 && page <= this.goodsTotalPages) {
                this.goodsCurrentPage = page;
            }
        },

        // ==================== 活动模板方法 ====================

        // 获取启用的商品列表（用于活动模板）
        async fetchGoodsForOrder() {
            try {
                const response = await axios.get('/api/goods/active-for-order', { withCredentials: true });
                this.availableGoodsForOrder = response.data.goods || [];
            } catch (err) {
                console.error('获取启用商品列表失败:', err);
            }
        },

        // 获取活动模板列表
        async fetchActivityTemplates() {
            this.loadingActivityTemplates = true;
            try {
                const response = await axios.get('/api/activity-templates', { withCredentials: true });
                this.activityTemplates = response.data.templates || [];
                this.filteredActivityTemplates = [...this.activityTemplates];
            } catch (err) {
                console.error('获取活动模板失败:', err);
            } finally {
                this.loadingActivityTemplates = false;
            }
        },

        // 获取启用状态的活动模板
        async fetchActiveActivityTemplates() {
            try {
                const response = await axios.get('/api/activity-templates/active', { withCredentials: true });
                this.activeActivityTemplates = response.data.templates || [];
            } catch (err) {
                console.error('获取启用活动模板失败:', err);
            }
        },

        // 搜索活动模板
        searchActivityTemplates() {
            this.filteredActivityTemplates = this.activityTemplates.filter(template => {
                const matchId = !this.activityTemplateFilters.id || template.id == this.activityTemplateFilters.id;
                const matchName = !this.activityTemplateFilters.name || template.name.includes(this.activityTemplateFilters.name);
                const matchType = !this.activityTemplateFilters.type || template.type == this.activityTemplateFilters.type;
                const matchStatus = this.activityTemplateFilters.status === '' || template.status == this.activityTemplateFilters.status;
                return matchId && matchName && matchType && matchStatus;
            });
            this.activityTemplateCurrentPage = 1;
        },

        // 重置活动模板筛选
        resetActivityTemplateFilters() {
            this.activityTemplateFilters = { id: '', name: '', type: '', status: '' };
            this.filteredActivityTemplates = [...this.activityTemplates];
            this.activityTemplateCurrentPage = 1;
        },

        // 活动模板分页
        changeActivityTemplatePage(page) {
            if (page >= 1 && page <= this.activityTemplateTotalPages) {
                this.activityTemplateCurrentPage = page;
            }
        },

        // 模板类型商品分页
        changeTemplateClassifyGoodsPage(page) {
            if (page >= 1 && page <= this.templateClassifyGoodsTotalPages) {
                this.templateClassifyGoodsCurrentPage = page;
            }
        },

        // 模板详情商品分页
        changeTemplateDetailGoodsPage(page) {
            if (page >= 1 && page <= this.templateDetailGoodsTotalPages) {
                this.templateDetailGoodsCurrentPage = page;
            }
        },

        // 打开活动模板详情抽屉
        async openActivityTemplateDetailDrawer(template) {
            try {
                const response = await axios.get(`/api/activity-templates/${template.id}`, { withCredentials: true });
                const data = response.data.template;

                let goods = [];
                if (data.select_type == 1) {
                    // 按类型选择，加载所有关联类型下的商品
                    const classifyList = data.classify_list || [];
                    for (const classify of classifyList) {
                        const goodsResponse = await axios.get(`/api/goods?classifyid=${classify.classify_id}&status=0`, { withCredentials: true });
                        goods = goods.concat(goodsResponse.data.goods || []);
                    }
                } else {
                    // 按商品选择
                    const goodsList = data.goods_list || [];
                    goods = goodsList.map(g => ({
                        id: g.goods_id,
                        name: g.goods_name,
                        price: g.price
                    }));
                }

                this.activityTemplateDetailData = {
                    id: data.id,
                    name: data.name,
                    type: data.type,
                    select_type: data.select_type,
                    status: data.status,
                    classifies: data.classify_list || [],
                    goods: goods
                };
                this.templateDetailGoodsCurrentPage = 1;
                this.showActivityTemplateDetailDrawer = true;
            } catch (err) {
                console.error('获取模板详情失败:', err);
                alert('获取模板详情失败');
            }
        },

        // 关闭活动模板详情抽屉
        closeActivityTemplateDetailDrawer() {
            this.showActivityTemplateDetailDrawer = false;
        },

        // 获取活动类型文本
        getActivityTypeText(type) {
            const types = { 1: '满减', 2: '满折', 3: '满赠', 4: '换购' };
            return types[type] || '';
        },

        // 打开新增活动模板抽屉
        async openAddActivityTemplateDrawer() {
            this.addActivityTemplateData = {
                name: '',
                type: '',
                select_type: '',
                status: 0,
                classify_ids: [],
                selected_classify_id: ''
            };
            this.selectedTemplateGoods = [];
            this.templateClassifyGoods = [];
            // 获取启用的类型列表
            try {
                const response = await axios.get('/api/classifies/active', { withCredentials: true });
                this.activeClassifies = response.data.classifies || [];
            } catch (err) {
                console.error('获取启用类型列表失败:', err);
            }
            this.fetchGoodsForOrder();
            this.showAddActivityTemplateDrawer = true;
        },

        // 关闭新增活动模板抽屉
        closeAddActivityTemplateDrawer() {
            this.showAddActivityTemplateDrawer = false;
        },

        // 选择方式变化
        onTemplateSelectTypeChange() {
            this.addActivityTemplateData.classify_ids = [];
            this.addActivityTemplateData.selected_classify_id = '';
            this.selectedTemplateGoods = [];
            this.templateClassifyGoods = [];
        },

        // 类型选择变化，加载类型下的商品
        async onTemplateClassifyChange() {
            const classifyId = this.addActivityTemplateData.selected_classify_id;
            if (!classifyId) {
                this.templateClassifyGoods = [];
                this.addActivityTemplateData.classify_ids = [];
                return;
            }
            // 设置classify_ids为选中的类型
            this.addActivityTemplateData.classify_ids = [classifyId];
            // 重置分页
            this.templateClassifyGoodsCurrentPage = 1;
            // 加载该类型下的商品
            try {
                const response = await axios.get(`/api/goods?classifyid=${classifyId}&status=0`, { withCredentials: true });
                this.templateClassifyGoods = response.data.goods || [];
            } catch (err) {
                console.error('获取类型商品失败:', err);
                this.templateClassifyGoods = [];
            }
        },

        // 编辑模板选择方式变化
        onEditTemplateSelectTypeChange() {
            this.editActivityTemplateData.classify_ids = [];
            this.editSelectedTemplateGoods = [];
        },

        // 打开添加模板商品弹窗
        openAddTemplateGoodsModal() {
            this.addTemplateGoodsData = {
                goods_id: '',
                name: '',
                brand_name: '',
                classify_name: '',
                price: 0
            };
            this.showAddTemplateGoodsModal = true;
        },

        // 关闭添加模板商品弹窗
        closeAddTemplateGoodsModal() {
            this.showAddTemplateGoodsModal = false;
        },

        // 模板商品选择变化
        onTemplateGoodsChange() {
            const selectedGoods = this.availableGoodsForOrder.find(g => g.id == this.addTemplateGoodsData.goods_id);
            if (selectedGoods) {
                this.addTemplateGoodsData.name = selectedGoods.name;
                this.addTemplateGoodsData.brand_name = selectedGoods.brand_name || '';
                this.addTemplateGoodsData.classify_name = selectedGoods.classify_name || '';
                this.addTemplateGoodsData.price = selectedGoods.price;
            }
        },

        // 保存模板商品
        saveTemplateGoods() {
            if (!this.addTemplateGoodsData.goods_id) {
                alert('请选择商品');
                return;
            }
            this.selectedTemplateGoods.push({
                goods_id: this.addTemplateGoodsData.goods_id,
                name: this.addTemplateGoodsData.name,
                brand_name: this.addTemplateGoodsData.brand_name,
                classify_name: this.addTemplateGoodsData.classify_name,
                price: this.addTemplateGoodsData.price
            });
            this.closeAddTemplateGoodsModal();
        },

        // 删除模板商品
        removeTemplateGoods(index) {
            this.selectedTemplateGoods.splice(index, 1);
        },

        // 保存新增活动模板
        async saveAddActivityTemplate() {
            if (!this.addActivityTemplateData.name || !this.addActivityTemplateData.type || !this.addActivityTemplateData.select_type) {
                alert('请填写必填项');
                return;
            }
            if (this.addActivityTemplateData.select_type == 1 && this.addActivityTemplateData.classify_ids.length === 0) {
                alert('请至少选择一个类型');
                return;
            }
            if (this.addActivityTemplateData.select_type == 2 && this.selectedTemplateGoods.length === 0) {
                alert('请至少选择一个商品');
                return;
            }
            try {
                const postData = {
                    name: this.addActivityTemplateData.name,
                    type: parseInt(this.addActivityTemplateData.type),
                    select_type: parseInt(this.addActivityTemplateData.select_type),
                    status: parseInt(this.addActivityTemplateData.status),
                    classify_ids: this.addActivityTemplateData.classify_ids,
                    goods_ids: this.selectedTemplateGoods.map(g => g.goods_id)
                };
                await axios.post('/api/activity-templates', postData, { withCredentials: true });
                await this.fetchActivityTemplates();
                this.closeAddActivityTemplateDrawer();
                alert('活动模板创建成功');
            } catch (err) {
                console.error('创建活动模板失败:', err);
                alert(err.response?.data?.error || '创建活动模板失败');
            }
        },

        // 打开编辑活动模板抽屉
        async openEditActivityTemplateDrawer(template) {
            try {
                const response = await axios.get(`/api/activity-templates/${template.id}`, { withCredentials: true });
                const data = response.data.template;
                this.editActivityTemplateData = {
                    id: data.id,
                    name: data.name,
                    type: data.type,
                    select_type: data.select_type,
                    status: data.status,
                    classify_ids: data.classify_list ? data.classify_list.map(c => c.classify_id) : []
                };
                if (data.goods_list) {
                    this.editSelectedTemplateGoods = data.goods_list.map(g => ({
                        goods_id: g.goods_id,
                        name: g.goods_name,
                        brand_name: g.brand_name || '',
                        classify_name: g.classify_name || '',
                        price: g.price
                    }));
                } else {
                    this.editSelectedTemplateGoods = [];
                }
                await this.fetchGoodsForOrder();
                this.showEditActivityTemplateDrawer = true;
            } catch (err) {
                console.error('获取活动模板详情失败:', err);
                alert('获取活动模板详情失败');
            }
        },

        // 关闭编辑活动模板抽屉
        closeEditActivityTemplateDrawer() {
            this.showEditActivityTemplateDrawer = false;
        },

        // 打开编辑模板商品弹窗
        openEditTemplateGoodsModal() {
            this.editTemplateGoodsData = {
                goods_id: '',
                name: '',
                brand_name: '',
                classify_name: '',
                price: 0
            };
            this.showEditTemplateGoodsModal = true;
        },

        // 关闭编辑模板商品弹窗
        closeEditTemplateGoodsModal() {
            this.showEditTemplateGoodsModal = false;
        },

        // 编辑模板商品选择变化
        onEditTemplateGoodsChange() {
            const selectedGoods = this.availableGoodsForOrder.find(g => g.id == this.editTemplateGoodsData.goods_id);
            if (selectedGoods) {
                this.editTemplateGoodsData.name = selectedGoods.name;
                this.editTemplateGoodsData.brand_name = selectedGoods.brand_name || '';
                this.editTemplateGoodsData.classify_name = selectedGoods.classify_name || '';
                this.editTemplateGoodsData.price = selectedGoods.price;
            }
        },

        // 保存编辑模板商品
        saveEditTemplateGoods() {
            if (!this.editTemplateGoodsData.goods_id) {
                alert('请选择商品');
                return;
            }
            this.editSelectedTemplateGoods.push({
                goods_id: this.editTemplateGoodsData.goods_id,
                name: this.editTemplateGoodsData.name,
                brand_name: this.editTemplateGoodsData.brand_name,
                classify_name: this.editTemplateGoodsData.classify_name,
                price: this.editTemplateGoodsData.price
            });
            this.closeEditTemplateGoodsModal();
        },

        // 删除编辑模板商品
        removeEditTemplateGoods(index) {
            this.editSelectedTemplateGoods.splice(index, 1);
        },

        // 保存编辑活动模板
        async saveEditActivityTemplate() {
            if (!this.editActivityTemplateData.name || !this.editActivityTemplateData.type || !this.editActivityTemplateData.select_type) {
                alert('请填写必填项');
                return;
            }
            if (this.editActivityTemplateData.select_type == 1 && this.editActivityTemplateData.classify_ids.length === 0) {
                alert('请至少选择一个类型');
                return;
            }
            if (this.editActivityTemplateData.select_type == 2 && this.editSelectedTemplateGoods.length === 0) {
                alert('请至少选择一个商品');
                return;
            }
            try {
                const putData = {
                    name: this.editActivityTemplateData.name,
                    type: parseInt(this.editActivityTemplateData.type),
                    select_type: parseInt(this.editActivityTemplateData.select_type),
                    status: parseInt(this.editActivityTemplateData.status),
                    classify_ids: this.editActivityTemplateData.classify_ids,
                    goods_ids: this.editSelectedTemplateGoods.map(g => g.goods_id)
                };
                await axios.put(`/api/activity-templates/${this.editActivityTemplateData.id}`, putData, { withCredentials: true });
                await this.fetchActivityTemplates();
                this.closeEditActivityTemplateDrawer();
                alert('活动模板更新成功');
            } catch (err) {
                console.error('更新活动模板失败:', err);
                alert(err.response?.data?.error || '更新活动模板失败');
            }
        },

        // 更新活动模板状态
        async updateActivityTemplateStatus(id, status) {
            try {
                await axios.put(`/api/activity-templates/${id}/status`, { status }, { withCredentials: true });
                await this.fetchActivityTemplates();
                alert(status === 0 ? '已启用' : '已禁用');
            } catch (err) {
                console.error('更新状态失败:', err);
                alert(err.response?.data?.error || '更新状态失败');
            }
        },

        // 删除活动模板
        async deleteActivityTemplate(id) {
            try {
                await axios.delete(`/api/activity-templates/${id}`, { withCredentials: true });
                await this.fetchActivityTemplates();
                alert('活动模板删除成功');
            } catch (err) {
                console.error('删除活动模板失败:', err);
                alert(err.response?.data?.error || '删除活动模板失败');
            }
        },

        // ==================== 活动管理方法 ====================

        // 获取活动列表
        async fetchActivities() {
            this.loadingActivities = true;
            try {
                const response = await axios.get('/api/activities', { withCredentials: true });
                this.activities = response.data.activities || [];
                this.filteredActivities = [...this.activities];
            } catch (err) {
                console.error('获取活动列表失败:', err);
            } finally {
                this.loadingActivities = false;
            }
        },

        // 搜索活动
        searchActivities() {
            this.filteredActivities = this.activities.filter(activity => {
                const matchId = !this.activityFilters.id || activity.id == this.activityFilters.id;
                const matchName = !this.activityFilters.name || activity.name.includes(this.activityFilters.name);
                const matchTemplate = !this.activityFilters.template_id || activity.template_id == this.activityFilters.template_id;
                const matchStatus = this.activityFilters.status === '' || activity.status == this.activityFilters.status;
                return matchId && matchName && matchTemplate && matchStatus;
            });
            this.activityCurrentPage = 1;
        },

        // 重置活动筛选
        resetActivityFilters() {
            this.activityFilters = { id: '', name: '', template_id: '', status: '' };
            this.filteredActivities = [...this.activities];
            this.activityCurrentPage = 1;
        },

        // 活动分页
        changeActivityPage(page) {
            if (page >= 1 && page <= this.activityTotalPages) {
                this.activityCurrentPage = page;
            }
        },

        // 格式化日期时间
        formatDateTime(dateStr) {
            if (!dateStr) return '';
            const date = new Date(dateStr);
            return date.toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit'
            });
        },

        // 打开新增活动抽屉
        async openAddActivityDrawer() {
            this.addActivityData = {
                name: '',
                template_id: '',
                template_type: '',
                template_select_type: '',
                start_time: '',
                end_time: '',
                status: 1,  // 默认禁用
                details: []
            };
            this.activityTemplateGoods = [];
            await this.fetchActiveActivityTemplates();
            this.showAddActivityDrawer = true;
        },

        // 关闭新增活动抽屉
        closeAddActivityDrawer() {
            this.showAddActivityDrawer = false;
        },

        // 活动模板选择变化
        async onActivityTemplateChange() {
            const template = this.activeActivityTemplates.find(t => t.id == this.addActivityData.template_id);
            if (template) {
                this.addActivityData.template_type = template.type;
                this.addActivityData.template_select_type = template.select_type;
                // 加载模板包含的商品
                await this.loadTemplateGoods(template.id, 'add');
            } else {
                this.addActivityData.template_type = '';
                this.addActivityData.template_select_type = '';
                this.activityTemplateGoods = [];
            }
        },

        // 加载模板包含的商品
        async loadTemplateGoods(templateId, mode) {
            try {
                const response = await axios.get(`/api/activity-templates/${templateId}`, { withCredentials: true });
                const template = response.data.template;
                let goods = [];

                if (template.select_type == 1) {
                    // 按类型选择，加载所有关联类型下的商品
                    const classifyList = template.classify_list || [];
                    for (const classify of classifyList) {
                        const goodsResponse = await axios.get(`/api/goods?classifyid=${classify.classify_id}&status=0`, { withCredentials: true });
                        goods = goods.concat(goodsResponse.data.goods || []);
                    }
                } else {
                    // 按商品选择，使用模板关联的商品
                    const goodsList = template.goods_list || [];
                    goods = goodsList.map(g => ({
                        id: g.goods_id,
                        name: g.goods_name,
                        price: g.price,
                        brand_name: g.brand_name,
                        classify_name: g.classify_name
                    }));
                }

                if (mode === 'add') {
                    this.activityTemplateGoods = goods;
                } else {
                    this.editActivityTemplateGoods = goods;
                }
            } catch (err) {
                console.error('获取模板商品失败:', err);
                if (mode === 'add') {
                    this.activityTemplateGoods = [];
                } else {
                    this.editActivityTemplateGoods = [];
                }
            }
        },

        // 加载订单活动模板包含的商品
        async loadOrderActivityTemplateGoods(activities) {
            try {
                let allGoods = [];
                const goodsMap = new Map(); // 用于去重

                for (const activity of activities) {
                    const templateId = activity.template_id;
                    if (!templateId) continue;

                    const response = await axios.get(`/api/activity-templates/${templateId}`, { withCredentials: true });
                    const template = response.data.template;
                    let goods = [];

                    if (template.select_type == 1) {
                        // 按类型选择，加载所有关联类型下的商品
                        const classifyList = template.classify_list || [];
                        for (const classify of classifyList) {
                            const goodsResponse = await axios.get(`/api/goods?classifyid=${classify.classify_id}&status=0`, { withCredentials: true });
                            goods = goods.concat(goodsResponse.data.goods || []);
                        }
                    } else {
                        // 按商品选择，使用模板关联的商品
                        const goodsList = template.goods_list || [];
                        goods = goodsList.map(g => ({
                            id: g.goods_id,
                            name: g.goods_name,
                            price: g.price,
                            brand_name: g.brand_name,
                            classify_name: g.classify_name
                        }));
                    }

                    // 添加活动名称到商品信息中，并去重
                    goods.forEach(g => {
                        if (!goodsMap.has(g.id)) {
                            goodsMap.set(g.id, {
                                ...g,
                                activities: [activity.name]
                            });
                        } else {
                            const existing = goodsMap.get(g.id);
                            if (!existing.activities.includes(activity.name)) {
                                existing.activities.push(activity.name);
                            }
                        }
                    });
                }

                // 转换Map为数组
                this.orderActivityTemplateGoods = Array.from(goodsMap.values());
            } catch (err) {
                console.error('获取订单活动模板商品失败:', err);
                this.orderActivityTemplateGoods = [];
            }
        },

        // 添加活动明细规则
        addActivityDetail() {
            this.addActivityData.details.push({ threshold_amount: 0, discount_value: 0 });
        },

        // 删除活动明细规则
        removeActivityDetail(index) {
            this.addActivityData.details.splice(index, 1);
        },

        // 保存新增活动
        async saveAddActivity() {
            if (!this.addActivityData.name || !this.addActivityData.template_id || !this.addActivityData.start_time || !this.addActivityData.end_time) {
                alert('请填写必填项');
                return;
            }
            if (this.addActivityData.template_type == 2 && this.addActivityData.details.length === 0) {
                alert('满折类型活动需要至少添加一条优惠规则');
                return;
            }
            try {
                const postData = {
                    name: this.addActivityData.name,
                    template_id: parseInt(this.addActivityData.template_id),
                    start_time: this.addActivityData.start_time,
                    end_time: this.addActivityData.end_time,
                    status: parseInt(this.addActivityData.status),
                    details: this.addActivityData.details
                };
                await axios.post('/api/activities', postData, { withCredentials: true });
                await this.fetchActivities();
                this.closeAddActivityDrawer();
                alert('活动创建成功');
            } catch (err) {
                console.error('创建活动失败:', err);
                alert(err.response?.data?.error || '创建活动失败');
            }
        },

        // 打开编辑活动抽屉
        async openEditActivityDrawer(activity) {
            try {
                await this.fetchActiveActivityTemplates();
                const response = await axios.get(`/api/activities/${activity.id}`, { withCredentials: true });
                const data = response.data.activity;

                // 格式化日期时间为input[type=datetime-local]格式（使用本地时间）
                const formatDateForInput = (dateStr) => {
                    if (!dateStr) return '';
                    const date = new Date(dateStr);
                    // 使用本地时间格式化，避免UTC时区转换问题
                    const year = date.getFullYear();
                    const month = String(date.getMonth() + 1).padStart(2, '0');
                    const day = String(date.getDate()).padStart(2, '0');
                    const hours = String(date.getHours()).padStart(2, '0');
                    const minutes = String(date.getMinutes()).padStart(2, '0');
                    return `${year}-${month}-${day}T${hours}:${minutes}`;
                };

                this.editActivityData = {
                    id: data.id,
                    name: data.name,
                    template_id: data.template_id,
                    template_type: data.template_type,
                    template_select_type: data.select_type,
                    start_time: formatDateForInput(data.start_time),
                    end_time: formatDateForInput(data.end_time),
                    status: data.status,
                    details: data.details || []
                };
                // 加载模板包含的商品
                if (data.template_id) {
                    await this.loadTemplateGoods(data.template_id, 'edit');
                }
                this.showEditActivityDrawer = true;
            } catch (err) {
                console.error('获取活动详情失败:', err);
                alert('获取活动详情失败');
            }
        },

        // 关闭编辑活动抽屉
        closeEditActivityDrawer() {
            this.showEditActivityDrawer = false;
        },

        // 打开活动详情抽屉
        async openActivityDetailDrawer(activity) {
            try {
                await this.fetchActiveActivityTemplates();
                const response = await axios.get(`/api/activities/${activity.id}`, { withCredentials: true });
                const data = response.data.activity;

                // 加载模板包含的商品
                let goods = [];
                if (data.template_id) {
                    const templateResponse = await axios.get(`/api/activity-templates/${data.template_id}`, { withCredentials: true });
                    const template = templateResponse.data.template;

                    if (template.select_type == 1) {
                        const classifyList = template.classify_list || [];
                        for (const classify of classifyList) {
                            const goodsResponse = await axios.get(`/api/goods?classifyid=${classify.classify_id}&status=0`, { withCredentials: true });
                            goods = goods.concat(goodsResponse.data.goods || []);
                        }
                    } else {
                        const goodsList = template.goods_list || [];
                        goods = goodsList.map(g => ({
                            id: g.goods_id,
                            name: g.goods_name,
                            price: g.price
                        }));
                    }
                }

                // 格式化日期时间
                const formatDateForDisplay = (dateStr) => {
                    if (!dateStr) return '';
                    return new Date(dateStr).toLocaleString('zh-CN');
                };

                this.activityDetailData = {
                    id: data.id,
                    name: data.name,
                    template_id: data.template_id,
                    template_name: data.template_name,
                    template_type: data.template_type,
                    template_select_type: data.select_type,
                    start_time: formatDateForDisplay(data.start_time),
                    end_time: formatDateForDisplay(data.end_time),
                    status: data.status,
                    details: data.details || [],
                    goods: goods
                };
                this.activityDetailGoodsCurrentPage = 1;
                this.showActivityDetailDrawer = true;
            } catch (err) {
                console.error('获取活动详情失败:', err);
                alert('获取活动详情失败');
            }
        },

        // 关闭活动详情抽屉
        closeActivityDetailDrawer() {
            this.showActivityDetailDrawer = false;
        },

        // 活动详情商品分页
        changeActivityDetailGoodsPage(page) {
            if (page >= 1 && page <= this.activityDetailGoodsTotalPages) {
                this.activityDetailGoodsCurrentPage = page;
            }
        },

        // 活动模板商品分页（新增/编辑活动）
        changeActivityTemplateGoodsPage(page) {
            if (page >= 1 && page <= this.activityTemplateGoodsTotalPages) {
                this.activityTemplateGoodsCurrentPage = page;
            }
        },

        // 编辑活动模板选择变化
        async onEditActivityTemplateChange() {
            const template = this.activeActivityTemplates.find(t => t.id == this.editActivityData.template_id);
            if (template) {
                this.editActivityData.template_type = template.type;
                this.editActivityData.template_select_type = template.select_type;
                // 加载模板包含的商品
                await this.loadTemplateGoods(template.id, 'edit');
            } else {
                this.editActivityData.template_type = '';
                this.editActivityData.template_select_type = '';
                this.editActivityTemplateGoods = [];
            }
        },

        // 添加编辑活动明细规则
        addEditActivityDetail() {
            this.editActivityData.details.push({ threshold_amount: 0, discount_value: 0 });
        },

        // 删除编辑活动明细规则
        removeEditActivityDetail(index) {
            this.editActivityData.details.splice(index, 1);
        },

        // 保存编辑活动
        async saveEditActivity() {
            if (!this.editActivityData.name || !this.editActivityData.template_id || !this.editActivityData.start_time || !this.editActivityData.end_time) {
                alert('请填写必填项');
                return;
            }
            if (this.editActivityData.template_type == 2 && this.editActivityData.details.length === 0) {
                alert('满折类型活动需要至少添加一条优惠规则');
                return;
            }
            try {
                const putData = {
                    name: this.editActivityData.name,
                    template_id: parseInt(this.editActivityData.template_id),
                    start_time: this.editActivityData.start_time,
                    end_time: this.editActivityData.end_time,
                    status: parseInt(this.editActivityData.status),
                    details: this.editActivityData.details
                };
                await axios.put(`/api/activities/${this.editActivityData.id}`, putData, { withCredentials: true });
                await this.fetchActivities();
                this.closeEditActivityDrawer();
                alert('活动更新成功');
            } catch (err) {
                console.error('更新活动失败:', err);
                alert(err.response?.data?.error || '更新活动失败');
            }
        },

        // 更新活动状态
        async updateActivityStatus(id, status) {
            try {
                await axios.put(`/api/activities/${id}/status`, { status }, { withCredentials: true });
                await this.fetchActivities();
                alert(status === 0 ? '已启用' : '已禁用');
            } catch (err) {
                console.error('更新状态失败:', err);
                alert(err.response?.data?.error || '更新状态失败');
            }
        },

        // 删除活动
        async deleteActivity(id) {
            try {
                await axios.delete(`/api/activities/${id}`, { withCredentials: true });
                await this.fetchActivities();
                alert('活动删除成功');
            } catch (err) {
                console.error('删除活动失败:', err);
                alert(err.response?.data?.error || '删除活动失败');
            }
        },

        // ==================== 合同管理方法 ====================

        // 获取合同列表
        async fetchContracts() {
            this.loadingContracts = true;
            try {
                const response = await axios.get('/api/contracts', { withCredentials: true });
                this.contracts = response.data.contracts || [];
                this.filteredContracts = [...this.contracts];
            } catch (err) {
                console.error('获取合同列表失败:', err);
            } finally {
                this.loadingContracts = false;
            }
        },

        // 搜索合同
        searchContracts() {
            this.filteredContracts = this.contracts.filter(contract => {
                const matchId = !this.contractFilters.id || contract.id == this.contractFilters.id;
                const matchStudentId = !this.contractFilters.student_id || contract.student_id == this.contractFilters.student_id;
                const matchStudentName = !this.contractFilters.student_name || (contract.student_name && contract.student_name.includes(this.contractFilters.student_name));
                const matchType = this.contractFilters.type === '' || contract.type == this.contractFilters.type;
                const matchStatus = this.contractFilters.status === '' || contract.status == this.contractFilters.status;
                const matchPaymentStatus = this.contractFilters.payment_status === '' || contract.payment_status == this.contractFilters.payment_status;

                return matchId && matchStudentId && matchStudentName && matchType && matchStatus && matchPaymentStatus;
            });
            this.contractCurrentPage = 1;
        },

        // 重置合同筛选
        resetContractFilters() {
            this.contractFilters = {
                id: '',
                student_id: '',
                student_name: '',
                type: '',
                status: '',
                payment_status: ''
            };
            this.filteredContracts = [...this.contracts];
            this.contractCurrentPage = 1;
        },

        // 合同分页切换
        changeContractPage(page) {
            if (page >= 1 && page <= this.contractTotalPages) {
                this.contractCurrentPage = page;
            }
        },

        // 获取合同类型文本
        getContractTypeText(type) {
            const typeMap = {
                0: '首报',
                1: '续报'
            };
            return typeMap[type] || '';
        },

        // 获取签署形式文本
        getSignatureFormText(form) {
            const formMap = {
                0: '线上签署',
                1: '线下签署'
            };
            return formMap[form] || '';
        },

        // 获取合同状态文本
        getContractStatusText(status) {
            const statusMap = {
                0: '待审核',
                50: '已通过',
                98: '已作废',
                99: '协议中止'
            };
            return statusMap[status] || '';
        },

        // 获取付款状态文本
        getPaymentStatusText(status) {
            const statusMap = {
                0: '未付款',
                10: '部分付款',
                30: '已付款'
            };
            return statusMap[status] || '';
        },

        // 打开新增合同抽屉
        async openAddContractDrawer() {
            this.addContractData = {
                student_id: '',
                student_name: '',
                type: '',
                signature_form: '',
                name: '',
                contract_amount: '',
                signatory: ''
            };

            // 获取启用的学生列表
            try {
                const response = await axios.get('/api/students/active', { withCredentials: true });
                this.activeStudentsForOrder = response.data.students;
            } catch (err) {
                console.error('获取学生列表失败:', err);
            }

            this.showAddContractDrawer = true;
        },

        // 关闭新增合同抽屉
        closeAddContractDrawer() {
            this.showAddContractDrawer = false;
            this.addContractData = {
                student_id: '',
                student_name: '',
                type: '',
                signature_form: '',
                name: '',
                contract_amount: '',
                signatory: ''
            };
        },

        // 学生选择变化（合同）
        onContractStudentChange() {
            const student = this.activeStudentsForOrder.find(s => s.id === parseInt(this.addContractData.student_id));
            if (student) {
                this.addContractData.student_name = student.student_name;
                this.updateContractName();
            }
        },

        // 更新合同名称
        updateContractName() {
            // 只有线下签署时才自动生成合同名称
            if (this.addContractData.signature_form == 1 && this.addContractData.student_id && this.addContractData.type !== '') {
                const typeText = this.addContractData.type == 0 ? '首报' : '续报';
                this.addContractData.name = `${this.addContractData.student_id}${this.addContractData.student_name}${typeText}合同`;
            } else if (this.addContractData.signature_form == 0) {
                // 线上签署时清空合同名称
                this.addContractData.name = '';
            }
        },

        // 确认提交合同
        confirmSubmitContract() {
            if (!this.addContractData.student_id || this.addContractData.type === '' ||
                this.addContractData.signature_form === '' || !this.addContractData.name ||
                !this.addContractData.contract_amount) {
                alert('请填写所有必填项');
                return;
            }

            // 线下签署需要填写签署方
            if (this.addContractData.signature_form == 1 && !this.addContractData.signatory) {
                alert('请填写签署方');
                return;
            }

            this.showSubmitContractConfirm = true;
        },

        // 提交合同
        async submitContract() {
            try {
                await axios.post('/api/contracts', this.addContractData, { withCredentials: true });
                alert('合同提交成功');
                this.showSubmitContractConfirm = false;
                this.closeAddContractDrawer();
                await this.fetchContracts();
            } catch (err) {
                console.error('提交合同失败:', err);
                alert(err.response?.data?.error || '提交合同失败');
            }
        },

        // 打开合同详情抽屉
        async openContractDetailDrawer(contract) {
            try {
                const response = await axios.get(`/api/contracts/${contract.id}`, { withCredentials: true });
                this.contractDetailData = response.data.contract;
                this.showContractDetailDrawer = true;
            } catch (err) {
                console.error('获取合同详情失败:', err);
                alert('获取合同详情失败');
            }
        },

        // 关闭合同详情抽屉
        closeContractDetailDrawer() {
            this.showContractDetailDrawer = false;
            this.contractDetailData = {
                id: '',
                student_id: '',
                student_name: '',
                type: '',
                signature_form: '',
                name: '',
                contract_amount: '',
                signatory: '',
                initiating_party: '',
                initiator: '',
                status: '',
                payment_status: '',
                create_time: ''
            };
        },

        // 确认撤销合同
        confirmRevokeContract() {
            this.showRevokeContractConfirm = true;
        },

        // 撤销合同
        async revokeContract() {
            try {
                await axios.put(`/api/contracts/${this.contractDetailData.id}/revoke`, {}, { withCredentials: true });
                alert('合同已撤销');
                this.showRevokeContractConfirm = false;
                this.closeContractDetailDrawer();
                await this.fetchContracts();
            } catch (err) {
                console.error('撤销合同失败:', err);
                alert(err.response?.data?.error || '撤销合同失败');
            }
        },

        // 打开中止合作弹窗
        openTerminateContractModal() {
            this.terminationAgreementFile = null;
            this.showTerminateContractModal = true;
        },

        // 关闭中止合作弹窗
        closeTerminateContractModal() {
            this.showTerminateContractModal = false;
            this.terminationAgreementFile = null;
        },

        // 处理中止协议文件上传
        handleTerminationFileUpload(event) {
            const file = event.target.files[0];
            if (file && file.type === 'application/pdf') {
                this.terminationAgreementFile = file;
            } else {
                alert('请上传PDF文件');
                event.target.value = '';
            }
        },

        // 中止合作
        async terminateContract() {
            if (!this.terminationAgreementFile) {
                alert('请上传中止协议');
                return;
            }

            try {
                // 这里简化处理，实际应该上传文件到服务器并获取文件路径
                // 暂时使用文件名作为路径
                const terminationAgreement = this.terminationAgreementFile.name;

                await axios.put(`/api/contracts/${this.contractDetailData.id}/terminate`, {
                    termination_agreement: terminationAgreement
                }, { withCredentials: true });

                alert('合作已中止');
                this.closeTerminateContractModal();
                this.closeContractDetailDrawer();
                await this.fetchContracts();
            } catch (err) {
                console.error('中止合作失败:', err);
                alert(err.response?.data?.error || '中止合作失败');
            }
        },

        // ==================== 收款管理方法 ====================

        // 获取收款列表
        async fetchPaymentCollections() {
            this.loadingPaymentCollections = true;
            try {
                const response = await axios.get('/api/payment-collections');
                this.paymentCollections = response.data.collections;
                this.filteredPaymentCollections = this.paymentCollections;
            } catch (err) {
                console.error('获取收款数据失败:', err);
            } finally {
                this.loadingPaymentCollections = false;
            }
        },

        // 搜索收款
        searchPaymentCollections() {
            this.filteredPaymentCollections = this.paymentCollections.filter(pc => {
                if (this.paymentCollectionFilters.id && pc.id != this.paymentCollectionFilters.id) {
                    return false;
                }
                if (this.paymentCollectionFilters.student_id && pc.student_id != this.paymentCollectionFilters.student_id) {
                    return false;
                }
                if (this.paymentCollectionFilters.order_id && pc.order_id != this.paymentCollectionFilters.order_id) {
                    return false;
                }
                if (this.paymentCollectionFilters.payer && pc.payer !== this.paymentCollectionFilters.payer) {
                    return false;
                }
                if (this.paymentCollectionFilters.payment_method !== '' && pc.payment_method != this.paymentCollectionFilters.payment_method) {
                    return false;
                }
                if (this.paymentCollectionFilters.trading_date) {
                    const tradingDate = pc.trading_hours ? pc.trading_hours.split('T')[0] : '';
                    if (tradingDate !== this.paymentCollectionFilters.trading_date) {
                        return false;
                    }
                }
                if (this.paymentCollectionFilters.status !== '' && pc.status != this.paymentCollectionFilters.status) {
                    return false;
                }
                return true;
            });
            this.paymentCollectionCurrentPage = 1;
        },

        // 重置收款筛选
        resetPaymentCollectionFilters() {
            this.paymentCollectionFilters = {
                id: '',
                student_id: '',
                order_id: '',
                payer: '',
                payment_method: '',
                trading_date: '',
                status: ''
            };
            this.filteredPaymentCollections = this.paymentCollections;
            this.paymentCollectionCurrentPage = 1;
        },

        // 分页切换
        changePaymentCollectionPage(page) {
            this.paymentCollectionCurrentPage = page;
        },

        // 获取付款场景文本
        getPaymentScenarioText(scenario) {
            const map = { 0: '线上', 1: '线下' };
            return map[scenario] || '-';
        },

        // 获取付款方式文本
        getPaymentMethodText(method) {
            const map = { 0: '微信', 1: '支付宝', 2: '优利支付', 3: '零零购支付', 9: '对公转账' };
            return map[method] || '-';
        },

        // 获取收款主体文本
        getPayeeEntityText(entity) {
            const map = { 0: '北京', 1: '西安' };
            return map[entity] || '-';
        },

        // 获取收款状态文本
        getCollectionStatusText(status) {
            const map = { 0: '待支付', 10: '未核验', 20: '已支付' };
            return map[status] || '-';
        },

        // 打开新增收款弹窗
        async openAddPaymentCollectionModal() {
            this.paymentCollectionForm = {
                student_id: '',
                order_id: '',
                pending_amount: '',
                payment_scenario: '',
                payment_method: '',
                payment_amount: '',
                payer: '',
                payee_entity: '',
                trading_hours: ''
            };
            this.studentUnpaidOrders = [];
            this.selectedOrderInfo = null;

            // 加载启用的学生列表
            try {
                const response = await axios.get('/api/students/active');
                this.activeStudents = response.data.students;
            } catch (err) {
                console.error('获取学生列表失败:', err);
                this.activeStudents = [];
            }

            this.showAddPaymentCollectionModal = true;
        },

        // 关闭新增收款弹窗
        closeAddPaymentCollectionModal() {
            this.showAddPaymentCollectionModal = false;
        },

        // 学生选择变化
        async onPaymentStudentChange() {
            this.paymentCollectionForm.order_id = '';
            this.paymentCollectionForm.pending_amount = '';
            this.selectedOrderInfo = null;

            if (!this.paymentCollectionForm.student_id) {
                this.studentUnpaidOrders = [];
                return;
            }

            try {
                const response = await axios.get(`/api/students/${this.paymentCollectionForm.student_id}/unpaid-orders`);
                this.studentUnpaidOrders = response.data.orders;
            } catch (err) {
                console.error('获取学生订单失败:', err);
                this.studentUnpaidOrders = [];
            }
        },

        // 订单选择变化
        async onPaymentOrderChange() {
            this.paymentCollectionForm.pending_amount = '';
            this.selectedOrderInfo = null;

            if (!this.paymentCollectionForm.order_id) {
                return;
            }

            try {
                const response = await axios.get(`/api/orders/${this.paymentCollectionForm.order_id}/pending-amount`);
                this.paymentCollectionForm.pending_amount = response.data.pending_amount;
                this.selectedOrderInfo = {
                    expected_payment_time: response.data.expected_payment_time
                };
            } catch (err) {
                console.error('获取待支付金额失败:', err);
            }
        },

        // 付款场景变化
        onPaymentScenarioChange() {
            this.paymentCollectionForm.payment_method = '';
        },

        // 提交新增收款
        async submitPaymentCollection() {
            // 校验必填项
            if (!this.paymentCollectionForm.student_id ||
                !this.paymentCollectionForm.order_id ||
                this.paymentCollectionForm.payment_scenario === '' ||
                this.paymentCollectionForm.payment_method === '' ||
                !this.paymentCollectionForm.payment_amount ||
                this.paymentCollectionForm.payee_entity === '') {
                alert('请填写所有必填项');
                return;
            }

            // 校验线下支付的商户订单号
            if (parseInt(this.paymentCollectionForm.payment_scenario) === 1) {
                if (!this.paymentCollectionForm.merchant_order) {
                    alert('请填写商户订单号');
                    return;
                }
            }

            // 校验付款金额
            const paymentAmount = parseFloat(this.paymentCollectionForm.payment_amount);
            const pendingAmount = parseFloat(this.paymentCollectionForm.pending_amount);
            if (paymentAmount > pendingAmount) {
                alert(`付款金额不能超过待支付金额(${pendingAmount})`);
                return;
            }

            // 校验线下支付的交易时间
            if (parseInt(this.paymentCollectionForm.payment_scenario) === 1) {
                if (!this.paymentCollectionForm.trading_hours) {
                    alert('请填写交易时间');
                    return;
                }

                // 检查交易时间是否与订单预计付款时间一致（只比较日期）
                if (this.selectedOrderInfo && this.selectedOrderInfo.expected_payment_time) {
                    const tradingDate = new Date(this.paymentCollectionForm.trading_hours);
                    const expectedDate = new Date(this.selectedOrderInfo.expected_payment_time);

                    // 只比较年月日
                    const tradingDateStr = tradingDate.toISOString().split('T')[0];
                    const expectedDateStr = expectedDate.toISOString().split('T')[0];

                    if (tradingDateStr !== expectedDateStr) {
                        alert('付款时间与订单不符！请重新填写！');
                        return;
                    }
                }
            }

            try {
                const payload = {
                    student_id: parseInt(this.paymentCollectionForm.student_id),
                    order_id: parseInt(this.paymentCollectionForm.order_id),
                    payment_scenario: parseInt(this.paymentCollectionForm.payment_scenario),
                    payment_method: parseInt(this.paymentCollectionForm.payment_method),
                    payment_amount: paymentAmount,
                    payer: this.paymentCollectionForm.payer,
                    payee_entity: parseInt(this.paymentCollectionForm.payee_entity),
                    merchant_order: this.paymentCollectionForm.merchant_order || null,
                    trading_hours: this.paymentCollectionForm.trading_hours || null
                };

                await axios.post('/api/payment-collections', payload);
                alert('收款新增成功');
                this.closeAddPaymentCollectionModal();
                await this.fetchPaymentCollections();
            } catch (err) {
                console.error('新增收款失败:', err);
                alert(err.response?.data?.error || '新增收款失败');
            }
        },

        // 确认到账
        confirmPaymentCollection(pc) {
            this.currentPaymentCollection = pc;
            this.showConfirmPaymentModal = true;
        },

        // 执行确认到账
        async doConfirmPayment() {
            try {
                await axios.put(`/api/payment-collections/${this.currentPaymentCollection.id}/confirm`);
                alert('已确认到账');
                this.showConfirmPaymentModal = false;
                this.currentPaymentCollection = null;
                await this.fetchPaymentCollections();
            } catch (err) {
                console.error('确认到账失败:', err);
                alert(err.response?.data?.error || '确认到账失败');
            }
        },

        // 确认删除收款
        confirmDeletePaymentCollection(pc) {
            this.currentPaymentCollection = pc;
            this.showDeletePaymentModal = true;
        },

        // 执行删除收款
        async doDeletePayment() {
            try {
                await axios.delete(`/api/payment-collections/${this.currentPaymentCollection.id}`);
                alert('收款已删除');
                this.showDeletePaymentModal = false;
                this.currentPaymentCollection = null;
                await this.fetchPaymentCollections();
            } catch (err) {
                console.error('删除收款失败:', err);
                alert(err.response?.data?.error || '删除收款失败');
            }
        },

        // ==================== 待认领收款管理方法 ====================

        // 获取待认领收款列表
        async fetchUnclaimedPayments() {
            this.loadingUnclaimedPayments = true;
            try {
                const params = new URLSearchParams();
                if (this.unclaimedFilters.id) params.append('id', this.unclaimedFilters.id);
                if (this.unclaimedFilters.payer) params.append('payer', this.unclaimedFilters.payer);
                if (this.unclaimedFilters.payment_method) params.append('payment_method', this.unclaimedFilters.payment_method);
                if (this.unclaimedFilters.arrival_date) params.append('arrival_date', this.unclaimedFilters.arrival_date);
                if (this.unclaimedFilters.claimer) params.append('claimer', this.unclaimedFilters.claimer);
                if (this.unclaimedFilters.status) params.append('status', this.unclaimedFilters.status);

                const response = await axios.get(`/api/unclaimed?${params.toString()}`);
                this.unclaimedPayments = response.data.unclaimed || [];
                this.filteredUnclaimedPayments = this.unclaimedPayments;
            } catch (err) {
                console.error('获取待认领列表失败:', err);
                alert(err.response?.data?.error || '获取待认领列表失败');
            } finally {
                this.loadingUnclaimedPayments = false;
            }
        },

        // 获取待认领付款方式文本
        getUnclaimedPaymentMethodText(method) {
            const methodMap = {
                0: '微信',
                1: '支付宝',
                2: '优利支付',
                3: '零零购支付',
                9: '对公转账'
            };
            return methodMap[method] || '未知';
        },

        // 获取待认领状态文本
        getUnclaimedStatusText(status) {
            const statusMap = {
                0: '待认领',
                1: '已认领'
            };
            return statusMap[status] || '未知';
        },

        // 确认认领
        confirmClaimUnclaimed(unclaimed) {
            this.currentUnclaimed = unclaimed;
            this.showClaimConfirmModal = true;
        },

        // 执行认领
        async doClaimUnclaimed() {
            // 验证订单ID是否填写
            if (!this.claimOrderId) {
                alert('请输入订单ID');
                return;
            }

            try {
                await axios.put(`/api/unclaimed/${this.currentUnclaimed.id}/claim`, {
                    order_id: parseInt(this.claimOrderId)
                });
                alert('认领成功');
                this.showClaimConfirmModal = false;
                this.currentUnclaimed = null;
                this.claimOrderId = '';
                await this.fetchUnclaimedPayments();
            } catch (err) {
                console.error('认领失败:', err);
                alert(err.response?.data?.error || '认领失败');
            }
        },

        // 确认删除待认领
        confirmDeleteUnclaimed(unclaimed) {
            this.currentUnclaimed = unclaimed;
            this.showDeleteUnclaimedModal = true;
        },

        // 执行删除待认领
        async doDeleteUnclaimed() {
            try {
                await axios.delete(`/api/unclaimed/${this.currentUnclaimed.id}`);
                alert('删除成功');
                this.showDeleteUnclaimedModal = false;
                this.currentUnclaimed = null;
                await this.fetchUnclaimedPayments();
            } catch (err) {
                console.error('删除失败:', err);
                alert(err.response?.data?.error || '删除失败');
            }
        },

        // 下载待认领Excel模板
        downloadUnclaimedTemplate() {
            window.location.href = '/api/unclaimed/template';
        },

        // 导入待认领Excel
        async importUnclaimedExcel(event) {
            const file = event.target.files[0];
            if (!file) return;

            const formData = new FormData();
            formData.append('file', file);

            try {
                const response = await axios.post('/api/unclaimed/import', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });

                if (response.data.errors && response.data.errors.length > 0) {
                    alert(response.data.message + '\n\n错误详情:\n' + response.data.errors.join('\n'));
                } else {
                    alert(response.data.message);
                }

                // 清空文件输入
                event.target.value = '';

                // 刷新列表
                await this.fetchUnclaimedPayments();
            } catch (err) {
                console.error('导入失败:', err);
                alert(err.response?.data?.error || '导入失败');
                event.target.value = '';
            }
        },

        // ============ 淘宝收款管理 ============
        // 获取淘宝已付款列表
        async fetchTaobaoPayments() {
            this.loadingTaobaoPayments = true;
            try {
                const params = new URLSearchParams();
                if (this.taobaoPaymentFilters.id) params.append('id', this.taobaoPaymentFilters.id);
                if (this.taobaoPaymentFilters.student_id) params.append('student_id', this.taobaoPaymentFilters.student_id);
                if (this.taobaoPaymentFilters.order_id) params.append('order_id', this.taobaoPaymentFilters.order_id);
                if (this.taobaoPaymentFilters.order_time) params.append('order_time', this.taobaoPaymentFilters.order_time);
                if (this.taobaoPaymentFilters.status !== '') params.append('status', this.taobaoPaymentFilters.status);

                const response = await axios.get(`/api/taobao-payments?${params.toString()}`);
                this.taobaoPayments = response.data.payments || [];
                this.filteredTaobaoPayments = this.taobaoPayments;
            } catch (err) {
                console.error('获取淘宝已付款列表失败:', err);
                alert(err.response?.data?.error || '获取淘宝已付款列表失败');
            } finally {
                this.loadingTaobaoPayments = false;
            }
        },

        // 搜索淘宝已付款
        searchTaobaoPayments() {
            this.filteredTaobaoPayments = this.taobaoPayments.filter(tp => {
                if (this.taobaoPaymentFilters.id && tp.id != this.taobaoPaymentFilters.id) {
                    return false;
                }
                if (this.taobaoPaymentFilters.student_id && tp.student_id != this.taobaoPaymentFilters.student_id) {
                    return false;
                }
                if (this.taobaoPaymentFilters.order_id && tp.order_id != this.taobaoPaymentFilters.order_id) {
                    return false;
                }
                if (this.taobaoPaymentFilters.order_time) {
                    const orderDate = tp.order_time ? tp.order_time.split('T')[0] : '';
                    if (orderDate !== this.taobaoPaymentFilters.order_time) {
                        return false;
                    }
                }
                if (this.taobaoPaymentFilters.status !== '' && tp.status != this.taobaoPaymentFilters.status) {
                    return false;
                }
                return true;
            });
            this.taobaoPaymentCurrentPage = 1;
        },

        // 重置淘宝已付款筛选
        resetTaobaoPaymentFilters() {
            this.taobaoPaymentFilters = {
                id: '',
                student_id: '',
                order_id: '',
                order_time: '',
                status: ''
            };
            this.fetchTaobaoPayments();
            this.taobaoPaymentCurrentPage = 1;
        },

        // 淘宝已付款分页切换
        changeTaobaoPaymentPage(page) {
            this.taobaoPaymentCurrentPage = page;
        },

        // 获取淘宝状态文本
        getTaobaoStatusText(status) {
            const map = { 0: '已下单', 10: '待认领', 20: '已认领', 30: '已到账', 40: '已退单' };
            return map[status] || '-';
        },

        // 打开新增淘宝收款弹窗
        async openAddTaobaoPaymentModal() {
            this.taobaoPaymentForm = {
                order_id: '',
                student_id: '',
                payer: '',
                zhifubao_account: '',
                payment_amount: '',
                order_time: '',
                merchant_order: '',
                pending_amount: ''
            };
            this.studentUnpaidOrders = [];
            this.selectedOrderInfo = null;

            // 加载启用的学生列表
            try {
                const response = await axios.get('/api/students/active');
                this.activeStudents = response.data.students;
            } catch (err) {
                console.error('获取学生列表失败:', err);
                this.activeStudents = [];
            }

            this.showAddTaobaoPaymentModal = true;
        },

        // 淘宝收款学生选择变化
        async onTaobaoStudentChange() {
            this.taobaoPaymentForm.order_id = '';
            this.taobaoPaymentForm.pending_amount = '';
            this.selectedOrderInfo = null;

            if (!this.taobaoPaymentForm.student_id) {
                this.studentUnpaidOrders = [];
                return;
            }

            try {
                const response = await axios.get(`/api/students/${this.taobaoPaymentForm.student_id}/unpaid-orders`);
                this.studentUnpaidOrders = response.data.orders;
            } catch (err) {
                console.error('获取学生订单失败:', err);
                this.studentUnpaidOrders = [];
            }
        },

        // 淘宝收款订单选择变化
        async onTaobaoOrderChange() {
            if (!this.taobaoPaymentForm.order_id) {
                this.taobaoPaymentForm.pending_amount = '';
                this.selectedOrderInfo = null;
                return;
            }

            const selectedOrder = this.studentUnpaidOrders.find(o => o.id == this.taobaoPaymentForm.order_id);
            if (selectedOrder) {
                this.selectedOrderInfo = selectedOrder;
                // 从后端API获取待支付金额（包括常规收款和淘宝收款）
                try {
                    const response = await axios.get(`/api/orders/${this.taobaoPaymentForm.order_id}/pending-amount`);
                    this.taobaoPaymentForm.pending_amount = response.data.pending_amount;
                } catch (err) {
                    console.error('获取待支付金额失败:', err);
                    alert('获取待支付金额失败');
                }
            }
        },

        // 新增淘宝收款
        async addTaobaoPayment() {
            // 验证必填字段
            if (!this.taobaoPaymentForm.order_id || !this.taobaoPaymentForm.student_id ||
                !this.taobaoPaymentForm.payment_amount || !this.taobaoPaymentForm.order_time) {
                alert('请填写所有必填字段');
                return;
            }

            try {
                await axios.post('/api/taobao-payments', this.taobaoPaymentForm);
                alert('新增成功');
                this.showAddTaobaoPaymentModal = false;
                await this.fetchTaobaoPayments();
            } catch (err) {
                console.error('新增失败:', err);
                alert(err.response?.data?.error || '新增失败');
            }
        },

        // 确认到账
        confirmTaobaoPayment(payment) {
            this.currentTaobaoPayment = payment;
            this.showConfirmTaobaoPaymentModal = true;
        },

        // 执行确认到账
        async doConfirmTaobaoPayment() {
            try {
                await axios.put(`/api/taobao-payments/${this.currentTaobaoPayment.id}/confirm`);
                alert('确认到账成功');
                this.showConfirmTaobaoPaymentModal = false;
                this.currentTaobaoPayment = null;
                await this.fetchTaobaoPayments();
            } catch (err) {
                console.error('确认到账失败:', err);
                alert(err.response?.data?.error || '确认到账失败');
            }
        },

        // 删除淘宝收款
        deleteTaobaoPayment(payment) {
            this.currentTaobaoPayment = payment;
            this.showDeleteTaobaoPaymentModal = true;
        },

        // 执行删除淘宝收款
        async doDeleteTaobaoPayment() {
            try {
                await axios.delete(`/api/taobao-payments/${this.currentTaobaoPayment.id}`);
                alert('删除成功');
                this.showDeleteTaobaoPaymentModal = false;
                this.currentTaobaoPayment = null;
                await this.fetchTaobaoPayments();
            } catch (err) {
                console.error('删除失败:', err);
                alert(err.response?.data?.error || '删除失败');
            }
        },

        // 获取淘宝待认领列表
        async fetchTaobaoUnclaimed() {
            this.loadingTaobaoUnclaimed = true;
            try {
                const params = new URLSearchParams();
                if (this.taobaoUnclaimedFilters.id) params.append('id', this.taobaoUnclaimedFilters.id);
                if (this.taobaoUnclaimedFilters.arrival_time) params.append('arrival_time', this.taobaoUnclaimedFilters.arrival_time);
                if (this.taobaoUnclaimedFilters.status !== '') params.append('status', this.taobaoUnclaimedFilters.status);

                const response = await axios.get(`/api/taobao-unclaimed?${params.toString()}`);
                this.taobaoUnclaimedPayments = response.data.unclaimed || [];
                this.filteredTaobaoUnclaimedPayments = this.taobaoUnclaimedPayments;
            } catch (err) {
                console.error('获取淘宝待认领列表失败:', err);
                alert(err.response?.data?.error || '获取淘宝待认领列表失败');
            } finally {
                this.loadingTaobaoUnclaimed = false;
            }
        },

        // 搜索淘宝待认领
        searchTaobaoUnclaimed() {
            this.filteredTaobaoUnclaimedPayments = this.taobaoUnclaimedPayments.filter(tu => {
                if (this.taobaoUnclaimedFilters.id && tu.id != this.taobaoUnclaimedFilters.id) {
                    return false;
                }
                if (this.taobaoUnclaimedFilters.arrival_time) {
                    const arrivalDate = tu.arrival_time ? tu.arrival_time.split('T')[0] : '';
                    if (arrivalDate !== this.taobaoUnclaimedFilters.arrival_time) {
                        return false;
                    }
                }
                if (this.taobaoUnclaimedFilters.status !== '' && tu.status != this.taobaoUnclaimedFilters.status) {
                    return false;
                }
                return true;
            });
            this.taobaoUnclaimedCurrentPage = 1;
        },

        // 重置淘宝待认领筛选
        resetTaobaoUnclaimedFilters() {
            this.taobaoUnclaimedFilters = {
                id: '',
                arrival_time: '',
                status: ''
            };
            this.fetchTaobaoUnclaimed();
            this.taobaoUnclaimedCurrentPage = 1;
        },

        // 淘宝待认领分页切换
        changeTaobaoUnclaimedPage(page) {
            this.taobaoUnclaimedCurrentPage = page;
        },

        // 获取淘宝待认领状态文本
        getTaobaoUnclaimedStatusText(status) {
            const map = { 10: '待认领', 20: '已认领' };
            return map[status] || '-';
        },

        // 认领淘宝待认领
        claimTaobaoUnclaimed(unclaimed) {
            this.currentTaobaoUnclaimed = unclaimed;
            this.showClaimTaobaoModal = true;
        },

        // 执行认领淘宝待认领
        async doClaimTaobaoUnclaimed() {
            if (!this.claimTaobaoOrderId) {
                alert('请输入订单ID');
                return;
            }

            try {
                await axios.put(`/api/taobao-unclaimed/${this.currentTaobaoUnclaimed.id}/claim`, {
                    order_id: parseInt(this.claimTaobaoOrderId)
                });
                alert('认领成功');
                this.showClaimTaobaoModal = false;
                this.currentTaobaoUnclaimed = null;
                this.claimTaobaoOrderId = '';
                await this.fetchTaobaoUnclaimed();
            } catch (err) {
                console.error('认领失败:', err);
                alert(err.response?.data?.error || '认领失败');
            }
        },

        // 删除淘宝待认领
        deleteTaobaoUnclaimed(unclaimed) {
            this.currentTaobaoUnclaimed = unclaimed;
            this.showDeleteTaobaoUnclaimedModal = true;
        },

        // 执行删除淘宝待认领
        async doDeleteTaobaoUnclaimed() {
            try {
                await axios.delete(`/api/taobao-unclaimed/${this.currentTaobaoUnclaimed.id}`);
                alert('删除成功');
                this.showDeleteTaobaoUnclaimedModal = false;
                this.currentTaobaoUnclaimed = null;
                await this.fetchTaobaoUnclaimed();
            } catch (err) {
                console.error('删除失败:', err);
                alert(err.response?.data?.error || '删除失败');
            }
        },

        // 下载淘宝待认领Excel模板
        downloadTaobaoTemplate() {
            window.location.href = '/api/taobao-unclaimed/template';
        },

        // 导入淘宝待认领Excel
        async importTaobaoExcel(event) {
            const file = event.target.files[0];
            if (!file) return;

            const formData = new FormData();
            formData.append('file', file);

            try {
                const response = await axios.post('/api/taobao-unclaimed/import', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });

                if (response.data.errors && response.data.errors.length > 0) {
                    alert(response.data.message + '\n\n错误详情:\n' + response.data.errors.join('\n'));
                } else {
                    alert(response.data.message);
                }

                // 清空文件输入
                event.target.value = '';

                // 刷新列表
                await this.fetchTaobaoUnclaimed();
            } catch (err) {
                console.error('导入失败:', err);
                alert(err.response?.data?.error || '导入失败');
                event.target.value = '';
            }
        },

        // ============ 分账明细管理 ============
        // 获取分账明细列表
        async fetchSeparateAccounts() {
            this.loadingSeparateAccounts = true;
            try {
                const params = new URLSearchParams();
                if (this.separateAccountFilters.id) params.append('id', this.separateAccountFilters.id);
                if (this.separateAccountFilters.uid) params.append('uid', this.separateAccountFilters.uid);
                if (this.separateAccountFilters.orders_id) params.append('orders_id', this.separateAccountFilters.orders_id);
                if (this.separateAccountFilters.childorders_id) params.append('childorders_id', this.separateAccountFilters.childorders_id);
                if (this.separateAccountFilters.goods_id) params.append('goods_id', this.separateAccountFilters.goods_id);
                if (this.separateAccountFilters.payment_id) params.append('payment_id', this.separateAccountFilters.payment_id);
                if (this.separateAccountFilters.payment_type !== '') params.append('payment_type', this.separateAccountFilters.payment_type);
                if (this.separateAccountFilters.type !== '') params.append('type', this.separateAccountFilters.type);

                const response = await axios.get(`/api/separate-accounts?${params.toString()}`);
                this.separateAccounts = response.data.separate_accounts || [];
                this.filteredSeparateAccounts = this.separateAccounts;
            } catch (err) {
                console.error('获取分账明细失败:', err);
                alert(err.response?.data?.error || '获取分账明细失败');
            } finally {
                this.loadingSeparateAccounts = false;
            }
        },

        // 搜索分账明细
        searchSeparateAccounts() {
            this.fetchSeparateAccounts();
            this.separateAccountCurrentPage = 1;
        },

        // 重置分账明细筛选
        resetSeparateAccountFilters() {
            this.separateAccountFilters = {
                id: '',
                uid: '',
                orders_id: '',
                childorders_id: '',
                goods_id: '',
                payment_id: '',
                payment_type: '',
                type: ''
            };
            this.fetchSeparateAccounts();
            this.separateAccountCurrentPage = 1;
        },

        // 分账明细分页切换
        changeSeparateAccountPage(page) {
            this.separateAccountCurrentPage = page;
        },

        // 获取分账类型文本
        getSeparateAccountTypeText(type) {
            const map = { 0: '售卖', 1: '冲回', 2: '退费' };
            return map[type] || '-';
        },

        // 获取收款类型文本
        getPaymentTypeText(paymentType) {
            const map = { 0: '常规收款', 1: '淘宝收款' };
            return map[paymentType] || '-';
        },

        // ==================== 角色管理相关方法 ====================

        // 获取角色列表
        async fetchRoles() {
            this.loadingRoles = true;
            try {
                const params = new URLSearchParams();
                if (this.roleFilters.id) params.append('id', this.roleFilters.id);
                if (this.roleFilters.name) params.append('name', this.roleFilters.name);
                if (this.roleFilters.status !== '') params.append('status', this.roleFilters.status);

                const response = await fetch(`/api/roles?${params.toString()}`, {
                    credentials: 'include'
                });
                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '获取角色列表失败');
                }
                const data = await response.json();
                this.roles = data.roles || [];
                this.filteredRoles = this.roles;
            } catch (error) {
                this.error = error.message;
                alert('获取角色列表失败：' + error.message);
            } finally {
                this.loadingRoles = false;
            }
        },

        // 筛选角色
        filterRoles() {
            this.fetchRoles();
            this.roleCurrentPage = 1;
        },

        // 重置角色筛选
        resetRoleFilters() {
            this.roleFilters = {
                id: '',
                name: '',
                status: ''
            };
            this.fetchRoles();
            this.roleCurrentPage = 1;
        },

        // 打开新增角色弹窗
        openAddRoleModal() {
            this.addRoleData = {
                name: '',
                comment: ''
            };
            this.showAddRoleModal = true;
        },

        // 关闭新增角色弹窗
        closeAddRoleModal() {
            this.showAddRoleModal = false;
        },

        // 提交新增角色
        async submitAddRole() {
            if (!this.addRoleData.name) {
                alert('请输入角色名称');
                return;
            }

            try {
                const response = await fetch('/api/roles', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include',
                    body: JSON.stringify(this.addRoleData)
                });

                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '创建角色失败');
                }

                alert('角色创建成功');
                this.closeAddRoleModal();
                this.fetchRoles();
            } catch (error) {
                alert('创建角色失败：' + error.message);
            }
        },

        // 打开编辑角色弹窗
        openEditRoleModal(role) {
            this.editRoleData = {
                id: role.id,
                name: role.name,
                comment: role.comment || ''
            };
            this.currentRole = role;
            this.showEditRoleModal = true;
        },

        // 关闭编辑角色弹窗
        closeEditRoleModal() {
            this.showEditRoleModal = false;
            this.currentRole = null;
        },

        // 提交编辑角色
        async submitEditRole() {
            if (!this.editRoleData.name) {
                alert('请输入角色名称');
                return;
            }

            try {
                const response = await fetch(`/api/roles/${this.editRoleData.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include',
                    body: JSON.stringify({
                        name: this.editRoleData.name,
                        comment: this.editRoleData.comment
                    })
                });

                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '更新角色失败');
                }

                alert('角色更新成功');
                this.closeEditRoleModal();
                this.fetchRoles();
            } catch (error) {
                alert('更新角色失败：' + error.message);
            }
        },

        // 打开启用角色确认
        openEnableRoleConfirm(role) {
            this.currentRole = role;
            this.showEnableRoleConfirm = true;
        },

        // 关闭启用角色确认
        closeEnableRoleConfirm() {
            this.showEnableRoleConfirm = false;
            this.currentRole = null;
        },

        // 确认启用角色
        async confirmEnableRole() {
            try {
                const response = await fetch(`/api/roles/${this.currentRole.id}/status`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include',
                    body: JSON.stringify({ status: 0 })
                });

                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '启用角色失败');
                }

                alert('角色已启用');
                this.closeEnableRoleConfirm();
                this.fetchRoles();
            } catch (error) {
                alert('启用角色失败：' + error.message);
            }
        },

        // 打开禁用角色确认
        openDisableRoleConfirm(role) {
            this.currentRole = role;
            this.showDisableRoleConfirm = true;
        },

        // 关闭禁用角色确认
        closeDisableRoleConfirm() {
            this.showDisableRoleConfirm = false;
            this.currentRole = null;
        },

        // 确认禁用角色
        async confirmDisableRole() {
            try {
                const response = await fetch(`/api/roles/${this.currentRole.id}/status`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include',
                    body: JSON.stringify({ status: 1 })
                });

                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '禁用角色失败');
                }

                alert('角色已禁用');
                this.closeDisableRoleConfirm();
                this.fetchRoles();
            } catch (error) {
                alert('禁用角色失败：' + error.message);
            }
        },

        // 打开权限配置弹窗
        async openRolePermissionsModal(role) {
            this.currentRole = role;
            this.showRolePermissionsModal = true;

            try {
                // 获取权限树
                const treeResponse = await fetch('/api/permissions/tree', {
                    credentials: 'include'
                });
                if (!treeResponse.ok) {
                    throw new Error('获取权限树失败');
                }
                const treeData = await treeResponse.json();
                this.permissionsTree = treeData.tree || [];

                // 构建可用权限列表（扁平化）
                this.availablePermissions = [];
                treeData.tree.forEach(menu => {
                    if (menu.permissions && menu.permissions.length > 0) {
                        menu.permissions.forEach(perm => {
                            this.availablePermissions.push({
                                id: perm.id,
                                name: perm.name,
                                menu_name: menu.name
                            });
                        });
                    }
                });

                // 获取角色已有权限
                const permResponse = await fetch(`/api/roles/${role.id}/permissions`, {
                    credentials: 'include'
                });
                if (!permResponse.ok) {
                    throw new Error('获取角色权限失败');
                }
                const permData = await permResponse.json();
                this.selectedRolePermissions = permData.permission_ids || [];

            } catch (error) {
                alert('加载权限数据失败：' + error.message);
                this.closeRolePermissionsModal();
            }
        },

        // 关闭权限配置弹窗
        closeRolePermissionsModal() {
            this.showRolePermissionsModal = false;
            this.currentRole = null;
            this.permissionsTree = [];
            this.selectedRolePermissions = [];
            this.availablePermissions = [];
        },

        // 选择权限（从可用列表移到已选列表）
        selectPermission(permissionId) {
            if (!this.selectedRolePermissions.includes(permissionId)) {
                this.selectedRolePermissions.push(permissionId);
            }
        },

        // 取消选择权限（从已选列表移除）
        deselectPermission(permissionId) {
            const index = this.selectedRolePermissions.indexOf(permissionId);
            if (index > -1) {
                this.selectedRolePermissions.splice(index, 1);
            }
        },

        // 获取可用权限列表（未被选中的）
        getAvailablePermissionsList() {
            return this.availablePermissions.filter(perm =>
                !this.selectedRolePermissions.includes(perm.id)
            );
        },

        // 获取已选权限列表
        getSelectedPermissionsList() {
            return this.availablePermissions.filter(perm =>
                this.selectedRolePermissions.includes(perm.id)
            );
        },

        // 提交权限配置
        async submitRolePermissions() {
            try {
                const response = await fetch(`/api/roles/${this.currentRole.id}/permissions`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include',
                    body: JSON.stringify({
                        permission_ids: this.selectedRolePermissions
                    })
                });

                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '更新权限失败');
                }

                alert('权限更新成功');
                this.closeRolePermissionsModal();
            } catch (error) {
                alert('更新权限失败：' + error.message);
            }
        },

        // 角色管理分页切换
        changeRolePage(page) {
            this.roleCurrentPage = page;
        },

        // 获取角色状态文本
        getRoleStatusText(status) {
            return status === 0 ? '启用' : '禁用';
        },

        // ==================== 菜单管理相关方法 ====================

        // 获取菜单管理列表
        async fetchMenuManagement() {
            this.loadingMenuManagement = true;
            try {
                const params = new URLSearchParams();
                if (this.menuManagementFilters.id) params.append('id', this.menuManagementFilters.id);
                if (this.menuManagementFilters.name) params.append('name', this.menuManagementFilters.name);
                if (this.menuManagementFilters.status !== '') params.append('status', this.menuManagementFilters.status);

                const response = await fetch(`/api/menu-management?${params.toString()}`, {
                    credentials: 'include'
                });
                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '获取菜单列表失败');
                }
                const data = await response.json();
                this.menuManagementList = data.menus || [];
                this.filteredMenuManagementList = this.menuManagementList;
            } catch (error) {
                this.error = error.message;
                alert('获取菜单列表失败：' + error.message);
            } finally {
                this.loadingMenuManagement = false;
            }
        },

        // 筛选菜单
        filterMenuManagement() {
            this.fetchMenuManagement();
            this.menuManagementCurrentPage = 1;
        },

        // 重置菜单筛选
        resetMenuManagementFilters() {
            this.menuManagementFilters = {
                id: '',
                name: '',
                status: ''
            };
            this.fetchMenuManagement();
            this.menuManagementCurrentPage = 1;
        },

        // 打开编辑菜单弹窗
        openEditMenuModal(menu) {
            this.editMenuData = {
                id: menu.id,
                name: menu.name,
                sort_order: menu.sort_order
            };
            this.currentMenu = menu;
            this.showEditMenuModal = true;
        },

        // 关闭编辑菜单弹窗
        closeEditMenuModal() {
            this.showEditMenuModal = false;
            this.currentMenu = null;
        },

        // 提交编辑菜单
        async submitEditMenu() {
            if (!this.editMenuData.name) {
                alert('请输入菜单名称');
                return;
            }

            if (this.editMenuData.sort_order === '') {
                alert('请输入排序');
                return;
            }

            try {
                const response = await fetch(`/api/menu-management/${this.editMenuData.id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include',
                    body: JSON.stringify({
                        name: this.editMenuData.name,
                        sort_order: parseInt(this.editMenuData.sort_order)
                    })
                });

                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '更新菜单失败');
                }

                alert('菜单更新成功');
                this.closeEditMenuModal();
                this.fetchMenuManagement();
            } catch (error) {
                alert('更新菜单失败：' + error.message);
            }
        },

        // 打开启用菜单确认
        openEnableMenuConfirm(menu) {
            this.currentMenu = menu;
            this.showEnableMenuConfirm = true;
        },

        // 关闭启用菜单确认
        closeEnableMenuConfirm() {
            this.showEnableMenuConfirm = false;
            this.currentMenu = null;
        },

        // 确认启用菜单
        async confirmEnableMenu() {
            try {
                const response = await fetch(`/api/menu-management/${this.currentMenu.id}/status`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include',
                    body: JSON.stringify({ status: 0 })
                });

                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '启用菜单失败');
                }

                alert('菜单已启用');
                this.closeEnableMenuConfirm();
                this.fetchMenuManagement();
            } catch (error) {
                alert('启用菜单失败：' + error.message);
            }
        },

        // 打开禁用菜单确认
        openDisableMenuConfirm(menu) {
            this.currentMenu = menu;
            this.showDisableMenuConfirm = true;
        },

        // 关闭禁用菜单确认
        closeDisableMenuConfirm() {
            this.showDisableMenuConfirm = false;
            this.currentMenu = null;
        },

        // 确认禁用菜单
        async confirmDisableMenu() {
            try {
                const response = await fetch(`/api/menu-management/${this.currentMenu.id}/status`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include',
                    body: JSON.stringify({ status: 1 })
                });

                if (!response.ok) {
                    const data = await response.json();
                    throw new Error(data.error || '禁用菜单失败');
                }

                alert('菜单已禁用');
                this.closeDisableMenuConfirm();
                this.fetchMenuManagement();
            } catch (error) {
                alert('禁用菜单失败：' + error.message);
            }
        },

        // 菜单管理分页切换
        changeMenuManagementPage(page) {
            this.menuManagementCurrentPage = page;
        },

        // 获取菜单状态文本
        getMenuStatusText(status) {
            return status === 0 ? '启用' : '禁用';
        }
    }
}).mount('#app');
