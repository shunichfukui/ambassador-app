import axios from 'axios'
import React, { useEffect, useState } from 'react'
import { Redirect } from 'react-router-dom'
import { User } from '../models/user'
import { Menu, Nav } from './organisms'

const Layout = (props: any) => {
    const [redirect, setRedirect] = useState(false)
    const [user, setUser] = useState<User | null>(null);

    useEffect(() => {
        (
            async () => {
                try {
                    const { data } = await axios.get('user');
                    setUser(data);
                } catch (error) {
                    setRedirect(true)
                }
            }
        )()
    }, [])

    if (redirect) {
        return <Redirect to={'/login'} />
    }

    return (
        <div>
        <Nav user={user} />
        <div className="container-fluid">
            <div className="row">
            <Menu />
            <main className="col-md-9 ms-sm-auto col-lg-10 px-md-4">
                <div className="table-responsive">
                    {props.children}
                </div>
            </main>
            </div>
        </div>
        </div>
    )
}

export default Layout