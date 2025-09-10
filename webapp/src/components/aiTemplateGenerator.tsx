// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useState} from 'react'
import {useHistory} from 'react-router-dom'

import octoClient from '../octoClient'

const AITemplateGenerator = (): JSX.Element => {
    const history = useHistory()
    const [prompt, setPrompt] = useState('')

    const onGenerate = async () => {
        const board = await octoClient.generateTemplate(prompt)
        if (board) {
            history.push(`/team/${board.teamId}/${board.id}`)
        }
    }

    return (
        <div>
            <input
                type='text'
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                placeholder='Describe your template'
            />
            <button onClick={onGenerate}>{'Generate Template'}</button>
        </div>
    )
}

export default AITemplateGenerator

