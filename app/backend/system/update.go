/*
 *
 *  * Copyright (c) 2025 Bronya0 <tangssst@163.com>. All rights reserved.
 *  * Original source: https://github.com/Bronya0
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *    http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package system

import (
	"app/backend/common"
	"app/backend/types"
	"context"
	"github.com/go-resty/resty/v2"
	"strings"
)

type Update struct {
	ctx context.Context
}

func (obj *Update) Start(ctx context.Context) {
	obj.ctx = ctx
}
func (obj *Update) CheckUpdate() *types.Tag {
	client := resty.New()
	tag := &types.Tag{}
	resp, err := client.R().SetResult(tag).Get(common.UPDATE_URL)
	if err != nil || resp.StatusCode() != 200 {
		return nil
	}
	tag.TagName = strings.TrimSpace(tag.TagName)
	return tag
}
